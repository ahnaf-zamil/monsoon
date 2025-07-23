import type React from "react";
import { Sidebar } from "../components/Sidebar";
import { Chat } from "../components/Chat";
import { useContext, useEffect } from "react";
import { AuthContext } from "../context/AuthContext";
import { useWebSocket } from "../context/SocketContext";
import type { IInboxEntry, IWebSocketDispatch } from "../ws/types";
import { OPCODES } from "../ws/opcodes";

export const Home: React.FC = () => {
    const user = useContext(AuthContext);
    const socket = useWebSocket();

    useEffect(() => {
        if (user && socket) {
            // TODO: add state management and separate socket event handlers
            const handleMessage = (e: MessageEvent) => {
                const payload: IWebSocketDispatch<IInboxEntry[]> = JSON.parse(
                    e.data
                );

                if (payload.opcode == OPCODES.RoomSync) {
                    console.log(payload.data)
                }
            };

            socket.addEventListener("message", handleMessage);
        }
    }, []);

    return (
        <>
            <div className="flex bg-darkbg max-h-svh">
                <Sidebar />
                <Chat />
            </div>
        </>
    );
};
