package broker

import "encoding/json"

type Message struct {
	data []byte
}

type MessageJson struct {
	Action string
	Data   string
}

func NewMessageJsonFromBytes(data []byte) (messageJson MessageJson, err error) {
	err = json.Unmarshal(data, &messageJson)
	return
}

func GetTopicAndPayload(data []byte) (topic []byte, payload []byte) {
	i := 0
	for _, char := range data {
		if char == '|' {
			break
		}
		i++
	}
	topic = data[:i]
	payload = data[i+1:]
	return
}
