# Jarvis

Service for Google Home actions. Deployed via Heroku using docker. The service must be publicly accessible in order to accept webhook traffic from Google. The service uses basic auth (cofigurable) for Google webhooks.

## Raspberry PI

### Setup

Flash your Raspberry Pi SD card using `client-pi-bakery.xml` with PiBakery. This recipe will download and install the client on your raspberry Pi from [http://your_server_here/dist/client-latest](http://your_server_here/dist/client-latest).

You will need to customize behavior for your Google Home device.

### Testing

Raspberry Pi hardware functions are mocked on Darwin so help local testing. To test your GPIO setup under Linux on a Pi:

```bash
go run cmd/pi/pi.go
```

## API

Base URL: `http://your_server_here`

* `/v1/christmas_lights` - handles Christmas light requests. Expects parameter `state-lights` to be either `on` or `off`.

The server will read the following environment values and override any flags:

* `BASIC_AUTH_USER` - Defaults to `raspberry`
* `BASIC_AUTH_PASSWORD` - Defaults to `password`
* `PORT` - Required

## Development

This project was created using `heroku create` for docker containers. The `Makefile` contains some useful build commands:

* `make build` - builds server for Linux amd64 and client for Linux ARM.
* `make docker` - builds the docker container
* `make check` - runs `go vet` and `go fmt`
* `make deploy` - builds and deploys container to heroku

### Running

The server expects clients to register with a device ID uuid.

To run local client against production:

```bash
go run cmd/client/client.go -addr <your server here> -device b9ea64b1-a00d-479f-92e4-e0f0a9a88692
```

Local testing client & server using docker:

```bash
docker run -it -p 8080:8080 -e PORT=8080 jarvis:latest
go run cmd/client/client.go -addr `docker-machine ip`:8080 -device b9ea64b1-a00d-479f-92e4-e0f0a9a88692
```

To run both client and server natively:

```bash
PORT=8080 go run cmd/server/server.go
go run cmd/client/client.go -device 9ed686c4-123b-4a43-a9cf-32b5cee6679c -addr localhost:8080
```
