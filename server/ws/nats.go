package ws

import (
	"encoding/json"
	"log"

	"monsoon/api"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

// All these functions will be stubbed during testing
type INATSPublisher interface {
	InitMsgListener()
	SendMsgNATS(content any)
}

type NATSPublisher struct {
	W IWebSocketHandler
}

func (n *NATSPublisher) InitNATS(natsURL string) (*nats.Conn, error) {
	/* Initializes NATS connection and starts message listener */
	var err error
	if nc, err = nats.Connect(natsURL); err != nil {
		log.Println("Error connecting to NATS:", err)
	}

	n.InitMsgListener()

	return nc, err
}

func (n *NATSPublisher) InitMsgListener() {
	/* This function will receive messages fromNATS and dispatch to all sockets in rooms */
	_, err := nc.Subscribe("message", func(m *nats.Msg) {
		// log.Printf("Received NATS msg: %s\n", string(m.Data))

		// Unmarshals received JSON into model for further processing
		var msg api.MessageModel
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			panic(err)
		}

		// Now we will start dispatching it to all sockets
		sock_list := GetSocketsForRoom(msg.ConversationID)

		// Looping over all sockets in the room and forwarding the message to them
		for _, s := range sock_list {
			if s == nil {
				continue
			}

			if err := n.W.DispatchEvent(s, OpMessageCreate, msg); err != nil {
				// Closing connection in case of write error, client-side should reconnect
				defer s.WsConn.Close()
			}
		}
	})
	if err != nil {
		log.Println("NATS subscription error:", err)
	}
}

func (n *NATSPublisher) SendMsgNATS(content any) {
	/* Dispatch new message to NATS cluster */

	payload, err := json.Marshal(content)
	if err == nil {
		err = nc.Publish("message", payload)
		if err != nil {
			log.Println("NATS publish error:", err)
		}
	}
}
