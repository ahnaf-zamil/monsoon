import type { IWebSocketDispatch } from "./types";

export const dispatchWebSocketEvent = (socket: WebSocket, opCode: string, data: any) => {
    const payload: IWebSocketDispatch<any> = {
        opcode: opCode,
        data
    }
    socket.send(JSON.stringify(payload))
}