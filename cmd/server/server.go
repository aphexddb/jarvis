package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/aphexddb/jarvis"
	"github.com/gorilla/mux"
)

const (
	clientDownloadPath   = "/dist"
	envBasicAuthUser     = "BASIC_AUTH_USER"
	envBasicAuthPassword = "BASIC_AUTH_PASSWORD"
)

// flags
var (
	userFlag     = flag.String("user", "raspberry", "Basic Auth username")
	passwordFlag = flag.String("password", "password", "Basic Auth password")
)

var hs *jarvis.HomeService
var port = os.Getenv("PORT")
var basicAuthUser, basicAuthPassword string

// basicAuth performs a Basic Auth check
func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	var realm = "Please enter your username and password for this site"
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(basicAuthUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(basicAuthPassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}
		handler(w, r)
	}
}

// use provides a cleaner interface for chaining middleware for single routes.
// Middleware functions are simple HTTP handlers (w http.ResponseWriter, r *http.Request)
//
//  r.HandleFunc("/login", use(loginHandler, rateLimit, csrf))
//  r.HandleFunc("/form", use(formHandler, csrf))
//  r.HandleFunc("/about", aboutHandler)
//
// See https://gist.github.com/elithrar/7600878#comment-955958 for how to extend it to suit simple http.Handler's
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// notFoundHandler handles 404's
func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "application/json")

	log.Println("Resource not found", req.URL.String(), "dumping request:")
	dumpRequest(req)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Resource not found. Latest client can be downloaded from /dist/client-latest"}`))
}

// fileHandler handles downloading files
func fileHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	file := vars["file"]
	log.Println("handling file download:", file, "("+req.URL.String()+")")

	// TODO: set linux executable content type
	// w.Header().Add("Content-Type", "application/x-executable")

	fs := http.FileServer(http.Dir(clientDownloadPath))
	http.StripPrefix("/dist/", fs).ServeHTTP(w, req)
}

// logWebhook logs a webhook event
func logWebhook(webhook jarvis.GoogleHomeWebhookRequest) {
	log.Println("WEBHOOK:",
		webhook.QueryResult.Action,
		"ACTION:",
		webhook.QueryResult.Action,
		"PARAMETERS:",
		webhook.QueryResult.Parameters,
		"QUERY:",
		webhook.QueryResult.QueryText,
	)
}

// dumpRequest dumps a HTTP request to the console
func dumpRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
}

// getEnvWithDefault gets an environment value with fallback value
func getEnvWithDefault(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	// force setting of PORT
	if port == "" {
		log.Println("$PORT must be set")
		os.Exit(1)
	}

	// get basic auth data
	basicAuthUser = getEnvWithDefault(envBasicAuthUser, *userFlag)
	basicAuthPassword = getEnvWithDefault(envBasicAuthPassword, *passwordFlag)

	// start websocket hub
	hub := jarvis.NewHub()
	go hub.Run()

	// create service
	hs = jarvis.NewHomeService(hub)
	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, req *http.Request) {
		jarvis.ServeWs(hub, w, req)
	})

	r := mux.NewRouter()
	r.HandleFunc("/v1/christmas_lights", use(hs.ChristmasLightsHandler, basicAuth))
	r.HandleFunc("/dist/{file}", fileHandler)

	// 404 handler
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.Handle("/", r)
	log.Println("Listening on", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
