import { useEffect } from "react";
import { useInboxStore } from "../store/inbox";
import { OPCODES } from "../ws/opcodes";
import type { IWebSocketDispatch, IInboxEntry } from "../ws/types";

export const useInboxSocketHandler = (socket: WebSocket | null) => {
    const inboxState = useInboxStore();

    useEffect(() => {
        if (socket) {
            const handleMessage = (e: MessageEvent) => {
                const payload: IWebSocketDispatch<IInboxEntry[]> = JSON.parse(
                    e.data,
                );
                switch (payload.opcode) {
                    case OPCODES.RoomSync:
                        inboxState.syncConversations(payload.data);
                }
            };

            socket.addEventListener("message", handleMessage);
            return () => socket.removeEventListener("message", handleMessage);
        }
    }, [socket, inboxState]);
};
