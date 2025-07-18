
export interface IWebSocketDispatch<EventDataType> {
    opcode: string;
    data: EventDataType;
}

export interface IHeartbeatInit {
  interval: number;
  timeout: number;
}