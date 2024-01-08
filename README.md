# Assignment Requirements

Your task is to write an application which acts as a pub/sub message broker server and communicates via QUIC.
 
Server specifications:
* Accepts QUIC connections on 2 ports (Publisher port, Subscriber port)
* The server notifies publishers if a subscriber has connected
* If no subscribers are connected, the server must inform the publishers
* The server sends any messages received from publishers to all connected subscribers

# How to Run

## Server

Run the command with `go` and specify path to configuration file, and certificate files in the command line arguments:
```bash
go run server/cmd/main.go server/config/base.yaml secrets/server.crt secrets/server.key
```
or alternatively run with `make`:
```bash
make run-server
```

## Publisher Client

Run the command with `go` and specify the publisher port as a cmd argument:
```bash
go run client/publisher/cmd/main.go 8081
```
or alternatively run with `make`:
```bash
make run-publisher
```

Once the publisher client is running, use the console input to publish messages.

## Subscriber Client

Run the command with `go` and specify the subscriber port as a cmd argument:
```bash
go run client/subscriber/cmd/main.go 8080
```
or alternatively run with `make`:
```bash
make run-subscriber
```

# Tests

Execute the command to run unit tests:
```bash
go test ./...
```
or alternatively run with `make`:
```bash
make test
```

Mocks were generated using the `mockgen` tool.

# Implementation Details

This project contains three different applications: server, publisher client, and subscriber client. Their corresponding code can be found in `server`, `client/publisher`, and `client/subscriber` directories accordingly. Common code can be found in `lib` directory.

## Server

Server can be configured via a configuration yaml file (`server/config/base.yaml`) that is loaded up on start up based on given configuration path.

On start up the `Server` creates two `Listeners`, one for subscribers and the other for publishers. `Listeners` start accepting incoming connections on separate goroutines. When a `Listener` accepts a new incoming connection then it executes the callback function provided by the `Server` and passes the connection along. Then the `Server` opens either a bi-directional stream (`ReadWriteStream`) for publishers, or a uni-directional write stream (`WriteStream`) for subscribers. When the stream is successfully opened, the `Server` passes it to the communication controller (`CommsController`).

`CommsController` is responsible for orchestrating the communication between subscribers and publishers. When `CommsController` receives a new publisher or subscriber stream, it then wraps the stream with a `notifier`. `notifier` runs in a separate goroutine, has a separate message queue, and is responsible for sending messages only to the given subscriber or publisher stream.

When `CommsController` receives a new message from a publisher it then puts the message to its' own message queue. `CommsController` processes the message queue in a separate goroutine. When processing a new message from a publisher, the `CommsController` will pass that message to all subscribers' `notifiers` to send the message independently of one another.

`CommsController` informs newly connected publishers of the current subscriber count (and if there aren't any). `CommsController` informs publishers when a new subscriber connects. `CommsController` sends the messages it receives from publishers to all subscribers.

`CommsController` maintains active subscribers and publishers, and removes them when they disconnect.

## Publisher Client

On start up the publisher `Client` connects to the server and accepts a bi-directional stream (`ReadWriteStream`). The `Client` will print out any messages it receives to the console output. Alternatively, a custom message receiver can be set by calling `Client.SetMessageReceiver`. New messages can be published to the server via `Client.Publish`.

The publisher client application will read console input and send the entered text to publishers on return (enter).

Publisher client will automatically shut down when the server shuts down.

## Subscriber Client

On start up the subscriber `Client` connects to the server and accepts a uni-directional read stream (`ReadStream`). The `Client` will print out any messages it receives to the console output. Alternatively, a custom message receiver can be set by calling `Client.SetMessageReceiver`.

Subscriber client will automatically shut down when the server shuts down.
