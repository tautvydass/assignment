run-server:
	go run server/cmd/main.go server/config/base.yaml secrets/server.crt secrets/server.key

run-subscriber:
	go run client/subscriber/cmd/main.go 8080
	
run-publisher:
	go run client/publisher/cmd/main.go 8081

test:
	go test ./...
