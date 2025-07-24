import { useContext, useEffect } from "react";
import { AuthContext } from "../context/AuthContext";
import { useWebSocket } from "../context/SocketContext";
import { useInboxStore } from "../store/inbox";
import { OPCODES } from "../ws/opcodes";
import type { IWebSocketDispatch, IInboxEntry } from "../ws/types";

export const useInboxSocketHandler = () => {
    const inboxState = useInboxStore();
    const user = useContext(AuthContext);
    const socket = useWebSocket();

    useEffect(() => {
        if (user && socket) {
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
        }
    }, []);
};
