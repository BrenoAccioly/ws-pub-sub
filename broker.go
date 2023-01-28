package broker

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

type Message struct {
	data []byte
}

type Broker struct {
	pathPattern string
	port        string
	server      http.Server
	connections []websocket.Conn
}

func NewBroker() Broker {
	broker := Broker{}
	broker.pathPattern = "/ws"
	broker.port = ":5000"
	broker.connections = make([]websocket.Conn, 0)
	return broker
}

func (broker Broker) ConnectionsSize() int {
	return len(broker.connections)
}

// publish to all connections
func (broker Broker) broadcast(message Message) {
	for _, ws := range broker.connections {
		ws.Write(message.data)
	}
}

func (broker *Broker) handleConnection(ws *websocket.Conn) {
	broker.connections = append(broker.connections, *ws)
	broker.broadcast(Message{[]byte("New connection Received")})
}

func (broker *Broker) Run() error {
	fmt.Println("Started server...")

	serveMux := http.NewServeMux()
	serveMux.Handle(
		broker.pathPattern,
		websocket.Handler(broker.handleConnection),
	)

	broker.server = http.Server{Addr: broker.port, Handler: serveMux}

	err := broker.server.ListenAndServe()

	return err
}

func (broker *Broker) Stop() {
	broker.server.Close()

	for _, conn := range broker.connections {
		conn.Close()
	}

	fmt.Println("Closed server...")
}
