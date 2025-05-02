package ws

import (
	"encoding/json"
	"log"

	"ws_realtime_app/lib"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

// All these functions will be stubbed during testing
type INATSPublisher interface {
	InitMsgListener()
	SendMsgNATS(content any)
}

type NATSPublisher struct{}

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
		// TODO: Process message and dispatch to connected sockets and rooms
		log.Printf("Received NATS msg: %s\n", string(m.Data))

		// Unmarshals received JSON into model for further processing
		// TODO: In case further message verification and integrity needs to be checked, do it here
		var msg lib.MessageModel
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			panic(err)
		}

		// Now we will start dispatching it to all sockets
		sock_list := GetSocketsForRoom(msg.RoomID)

		// Looping over all sockets in the room and forwarding the message to them
		for _, s := range sock_list {
			if err := s.WsConn.WriteJSON(msg); err != nil {
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
