package broker

import (
	"testing"

	broker "github.com/BrenoAccioly/ws-pub-sub"
	"golang.org/x/net/websocket"
)

var origin string = "http://localhost/"
var url string = "ws://localhost:5000/ws"

func TestNewConnection(t *testing.T) {
	broker := broker.NewBroker()

	go broker.Run()
	defer broker.Stop()

	_, err := websocket.Dial(url, "", origin)

	if err != nil {
		t.Errorf("Connection Failed")
	}

	if broker.ConnectionsSize() != 1 {
		t.Error()
	}
}

func TestMessageBroadcast(t *testing.T) {
	broker := broker.NewBroker()

	go broker.Run()
	defer broker.Stop()

	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		t.Errorf("Connection Failed")
	}

	var msg = make([]byte, 512)
	//var n int
	if _, err = ws.Read(msg); err != nil {
		t.Errorf("Read")
	}
	//log.Fatalf("Received: %s.\n", msg[:n])
}
