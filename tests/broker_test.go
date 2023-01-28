package broker

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	bk "github.com/BrenoAccioly/ws-pub-sub"
	"golang.org/x/net/websocket"
)

var origin string = "http://localhost/"
var url string = "ws://localhost:5000/ws"

func TestNewConnection(t *testing.T) {
	broker := bk.NewBroker()

	go broker.Run()
	defer broker.Stop()

	time.Sleep(time.Second)
	_, err := websocket.Dial(url, "", origin)

	if err != nil {
		t.Errorf("Connection Failed: %s", err)
	}

	if broker.ConnectionsSize() != 1 {
		t.Error()
	}
}

func TestSubscription(t *testing.T) {
	broker := bk.NewBroker()

	go broker.Run()
	defer broker.Stop()

	time.Sleep(time.Second)
	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		t.Errorf("Connection Failed: %s", err)
	}

	var messageBuffer = make([]byte, 512)

	// subscribe

	messageJson := bk.MessageJson{
		Action: "subscribe", Data: "topic1",
	}

	messageJsonBytes, _ := json.Marshal(messageJson)
	ws.Write(messageJsonBytes)

	// subcription success

	if n, err := ws.Read(messageBuffer); err != nil {
		t.Errorf("Read")
	} else {
		fmt.Printf("[Client] %s", messageBuffer[:n])
	}
}

func TestPublish(t *testing.T) {
	broker := bk.NewBroker()

	go broker.Run()
	defer broker.Stop()

	time.Sleep(time.Second)
	ws1, err := websocket.Dial(url, "", origin)
	ws2, err := websocket.Dial(url, "", origin)

	if err != nil {
		t.Errorf("Connection Failed: %s", err)
	}

	var messageBuffer = make([]byte, 512)

	// subscribe

	messageJson := bk.MessageJson{
		Action: "subscribe", Data: "topic1",
	}

	messageJsonBytes, _ := json.Marshal(messageJson)
	ws1.Write(messageJsonBytes)

	// subcription success

	if n, err := ws1.Read(messageBuffer); err != nil {
		t.Errorf("Read")
	} else {
		fmt.Printf("[Client 1] %s\n", messageBuffer[:n])
	}

	// publish

	messageJson = bk.MessageJson{
		Action: "publish", Data: "topic1|hello",
	}

	messageJsonBytes, _ = json.Marshal(messageJson)
	ws2.Write(messageJsonBytes)

	if n, err := ws1.Read(messageBuffer); err != nil {
		t.Errorf("Read")
	} else {
		fmt.Printf("[Client 1] %s\n", messageBuffer[:n])
	}
}
