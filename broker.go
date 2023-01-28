package broker

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/slices"
	"golang.org/x/net/websocket"
)

type Message struct {
	data []byte
}

type MessageJson struct {
	Action string
	Data   string
}

type Connection struct {
	ID   int64
	Conn *websocket.Conn
}

func NewConnection(ws *websocket.Conn) Connection {
	var n int64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	conn := Connection{}
	conn.ID = n
	conn.Conn = ws
	return conn
}

type Broker struct {
	pathPattern string
	port        string
	server      http.Server
	connections map[int64]Connection
	topics      map[string][]int64
}

func NewBroker() Broker {
	broker := Broker{}
	broker.pathPattern = "/ws"
	broker.port = ":5000"
	broker.connections = make(map[int64]Connection)
	broker.topics = make(map[string][]int64)
	return broker
}

func (broker Broker) ConnectionsSize() int {
	return len(broker.connections)
}

// publish to all connections
func (broker Broker) broadcast(message Message) {
	for _, conn := range broker.connections {
		conn.Conn.Write(message.data)
	}
}

func (broker Broker) NewSubscription(connID int64, topic string) {
	fmt.Println("[Server] Received new subscription")
	if !slices.Contains(broker.topics[topic], connID) {
		broker.topics[topic] = append(broker.topics[topic], connID)
	}
	broker.connections[connID].Conn.Write([]byte("Subscription Success\n"))
}

func (broker Broker) publish(data []byte) {
	fmt.Println("[Server] Received publish")
	topic := ""
	i := 0
	for _, c := range data {
		if c == '|' {
			break
		}
		topic += string(c)
		i++
	}
	payload := data[i+1:]
	connIDs, ok := broker.topics[topic]
	if ok {
		for _, id := range connIDs {
			fmt.Printf("[Server] sent to %d\n", id)
			broker.connections[id].Conn.Write(payload)
		}
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

		var messageJson MessageJson
		err = json.Unmarshal(messageBuffer[:n], &messageJson)

		switch messageJson.Action {
		case "subscribe":
			broker.NewSubscription(conn.ID, messageJson.Data)
		case "publish":
			broker.publish([]byte(messageJson.Data))
		default:
		}
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
