import { useEffect } from "react";
import type { IWebSocketDispatch } from "../ws/types";
import type { IMessageData } from "../api/types";
import { OPCODES } from "../ws/opcodes";
import { useMessageStore } from "../store/message";

export const useMessageSocketHandler = (socket: WebSocket | null) => {
    const messageStore = useMessageStore();

    useEffect(() => {
        if (socket) {
            const handleMessage = (e: MessageEvent) => {
                const payload: IWebSocketDispatch<IMessageData> = JSON.parse(
                    e.data,
                );
                switch (payload.opcode) {
                    case OPCODES.MessageCreate:
                        messageStore.storeMessages(
                            payload.data.conversation_id,
                            [payload.data],
                        );
                }
            };

            socket.addEventListener("message", handleMessage);
        }
    }, [socket]);
};
