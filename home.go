package jarvis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// HomeService represents a home automation service
type HomeService struct {
	hub *Hub
}

// NewHomeService creates a new home service
func NewHomeService(hub *Hub) *HomeService {
	return &HomeService{
		hub: hub,
	}
}

// logWebhook logs a webhook event
func logWebhook(webhook GoogleHomeWebhookRequest) {
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

// ChristmasLightsHandler is the christmas light webhook
func (j *HomeService) ChristmasLightsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var webhook GoogleHomeWebhookRequest
	err := decoder.Decode(&webhook)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Println("Unable to read webhook")
		return
	}

	logWebhook(webhook)

	stateLights, ok := webhook.QueryResult.Parameters["state-lights"]
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(`{"message": "state-lights parameter expected"}`))
		log.Println("state-lights parameter missing")
		return
	}

	if stateLights != "on" && stateLights != "off" && stateLights != "" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(`{"message": "state-lights parameter value unknown"}`))
		log.Println("state-lights value unknown:", stateLights)
		return
	}

	// TODO : send message to client
	w.Write([]byte(`{"message": "Webhook handled successfully. TODO: SEND TO CLIENT"}`))
}
