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
	connections []websocket.Conn
}

func NewBroker() Broker {
	broker := Broker{}
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

func (broker *Broker) Run() {
	fmt.Println("Started server...")

	http.Handle("/ws", websocket.Handler(broker.handleConnection))

	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func (broker *Broker) Stop() {
	// TODO
}
