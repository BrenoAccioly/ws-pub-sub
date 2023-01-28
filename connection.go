package broker

import (
	"crypto/rand"
	"encoding/binary"

	"golang.org/x/net/websocket"
)

type UUID int64

type Connection struct {
	ID   UUID
	Conn *websocket.Conn
}

func NewConnection(ws *websocket.Conn) Connection {
	var n UUID
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	conn := Connection{}
	conn.ID = n
	conn.Conn = ws
	return conn
}
