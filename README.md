# Assignment Requirements

Your task is to write an application which acts as a pub/sub message broker server and communicates via QUIC.
 
Server specifications:
* Accepts QUIC connections on 2 ports (Publisher port, Subscriber port)
* The server notifies publishers if a subscriber has connected
* If no subscribers are connected, the server must inform the publishers
* The server sends any messages received from publishers to all connected subscribers

# How to Run

## Server

Execute the command:
```bash
go run server/cmd/main.go server/config/base.yaml
```

## Publisher Client

Execute the command:
```bash
go run client/publisher/cmd/main.go
```

## Subscriber Client

Execute the command:
```bash
go run client/subscriber/cmd/main.go
```

# Tests

Execute the command to run unit tests:
```bash
go test ./...
```