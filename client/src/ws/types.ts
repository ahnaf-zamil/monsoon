export interface IWebSocketDispatch<EventDataType> {
    opcode: string;
    data: EventDataType;
}

export interface IHeartbeatInit {
    interval: number;
    timeout: number;
}

type ConversationType = "GROUP" | "DM";

export interface IInboxEntry {
    conversation_id: string;
    name: string;
    updated_at: number;
    user_id: string | null;
    type: ConversationType;
}
