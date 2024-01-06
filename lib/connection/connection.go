package connection

import "github.com/quic-go/quic-go"

// Connection is an interface for the connection.
type Connection interface{}

type connection struct {
	conn quic.Connection
}

// New constructs a new connection.
func New(conn quic.Connection) Connection {
	return &connection{
		conn: conn,
	}
}
