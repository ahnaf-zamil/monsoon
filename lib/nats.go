package lib

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func InitNATS(natsURL string) *nats.Conn {
	/* Initializes NATS connection and starts message listener */
	var err error
	if nc, err = nats.Connect(natsURL); err != nil {
		log.Println("Error connecting to NATS:", err)
	}

	InitMsgListener()

	return nc
}

func InitMsgListener() {
	/* This function will receive messages fromNATS and dispatch to all sockets in rooms */
	nc.Subscribe("message", func(m *nats.Msg) {
		log.Printf("Received NATS msg: %s\n", string(m.Data))
	})
}

func SendMsgNATS(content any) {
	payload, err := json.Marshal(content)
	if err == nil {
		nc.Publish("message", payload)
	}
}
