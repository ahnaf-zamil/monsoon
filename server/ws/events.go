package ws

import (
	"log"
	"time"
)

func (w *WebSocketHandler) DispatchEvent(socket *Socket, opCode EventOpCode, data any) error {
	payload := &WebSocketEvent{OpCode: opCode, Data: data}
	if err := socket.WsConn.WriteJSON(payload); err != nil {
		defer socket.WsConn.Close()
		return err
	}
	return nil
}

func (w *WebSocketHandler) StartHeartbeat() {
	/* Cleans up dead connections who have not sent a "heartbeat" within timeout

	A socket's last heartbeat timeout is updated by sending heartbeat events every HEARTBEAT_INTERVAL_CLIENT seconds
	*/

	log.Println("Starting socket heartbeat checks")
	go func() {
		ticker := time.NewTicker(HEARTBEAT_CHECK_INTERVAL)
		defer ticker.Stop()

		for {
			<-ticker.C
			mu.RLock()
			for _, sock := range socketList {
				if time.Since(sock.LastHeartbeat) > HEARTBEAT_TIMEOUT {
					// If the socket has not made a heartbeat in the within timeout, disconnect it to save memory
					log.Println("Socket", sock.ID, "has been sent no heartbeat since the last", time.Since(sock.LastHeartbeat), ". Disconnecting...")
					sock.WsConn.Close()
					RemoveSocketFromList(sock)
				}
			}
			mu.RUnlock()
		}
	}()
}

func (w *WebSocketHandler) HandleClientHeartbeat(socket *Socket) {
	// Respond with heartbeat ack
	socket.LastHeartbeat = time.Now()
	_ = w.DispatchEvent(socket, OpHeartbeatAck, nil)
}
