import type { IWebSocketDispatch, IHeartbeatInit } from "../ws/types";
import { OPCODES } from "../ws/opcodes";
import { dispatchWebSocketEvent } from "../ws/events";

let isHeartbeating = false;
let heartbeatTimeoutId: ReturnType<typeof setTimeout> | null = null;

export const useWSHeartbeat = (socket: WebSocket | null) => {
  if (!socket) return;
  if (isHeartbeating) return;

  const onMessage = (e: MessageEvent) => {
    const payload: IWebSocketDispatch<IHeartbeatInit> = JSON.parse(e.data);
    if (payload.opcode == OPCODES.HeartbeatInit) {
      console.log(
        `Starting heartbeat with ${payload.data.interval}ms interval`
      );

      setInterval(() => {
        if (heartbeatTimeoutId) clearTimeout(heartbeatTimeoutId);
        heartbeatTimeoutId = setTimeout(() => {
          console.warn("No heartbeat ack received within timeout. Disconnecting...");
          isHeartbeating = false;
          socket.close(); // Implement reconnect later
        }, payload.data.timeout);

        if (!socket || socket.readyState !== WebSocket.OPEN) return;
        dispatchWebSocketEvent(socket, OPCODES.Heartbeat, null);
        console.log("Dispatched heartbeat");
      }, payload.data.interval);
      isHeartbeating = true;

    } else if (payload.opcode == OPCODES.HeartbeatAck) {
      console.log("Heartbeat ack received");
      if (heartbeatTimeoutId) clearTimeout(heartbeatTimeoutId);
    }
  };

  socket.addEventListener("message", onMessage);
};
