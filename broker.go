package broker

import (
	"fmt"
	"net/http"

	"golang.org/x/exp/slices"
	"golang.org/x/net/websocket"
)

type Broker struct {
	pathPattern string
	port        string
	server      http.Server
	connections map[UUID]Connection
	topics      map[string][]UUID
}

func NewBroker() Broker {
	broker := Broker{}
	broker.pathPattern = "/ws"
	broker.port = ":5000"
	broker.connections = make(map[UUID]Connection)
	broker.topics = make(map[string][]UUID)
	return broker
}

func (broker Broker) ConnectionsSize() int {
	return len(broker.connections)
}

func (broker Broker) broadcast(message Message) {
	for _, conn := range broker.connections {
		conn.Conn.Write(message.data)
	}
}

func (broker Broker) NewSubscription(connID UUID, topic string) {
	fmt.Println("[Server] Received new subscription")
	if !slices.Contains(broker.topics[topic], connID) {
		broker.topics[topic] = append(broker.topics[topic], connID)
	}
	broker.connections[connID].Conn.Write([]byte("Subscription Success\n"))
}

func (broker Broker) publish(data []byte) {
	fmt.Println("[Server] Received publish")

	topic, payload := GetTopicAndPayload(data)

	connIDs, ok := broker.topics[string(topic)]
	if ok {
		for _, id := range connIDs {
			fmt.Printf("[Server] sent to %d\n", id)
			broker.connections[id].Conn.Write(payload)
		}
	}
}

func (broker *Broker) handleMessageJson(id UUID, messageJson MessageJson) {
	switch messageJson.Action {
	case "subscribe":
		broker.NewSubscription(id, messageJson.Data)
	case "publish":
		broker.publish([]byte(messageJson.Data))
	case "unsubscribe":
	default:
	}
}

func (broker *Broker) handleConnection(ws *websocket.Conn) {
	conn := NewConnection(ws)
	broker.connections[conn.ID] = conn
	//broker.broadcast(Message{[]byte("New connection Received")})

	for {
		messageBuffer := make([]byte, 512)
		n, err := ws.Read(messageBuffer)

		if err != nil {
			return
		}

		fmt.Printf("[Server] received: %s\n", messageBuffer)

		messageJson, err := NewMessageJsonFromBytes(messageBuffer[:n])

		broker.handleMessageJson(conn.ID, messageJson)
	}
}

func (broker *Broker) Run() error {
	fmt.Println("[Server] Started server...")

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
		conn.Conn.Close()
	}

	fmt.Println("[Server] Closed server...")
}
