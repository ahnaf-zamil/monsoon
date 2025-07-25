import {
    createContext,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from "react";
import { useCurrentUser } from "./AuthContext";
import { getAccessToken } from "../api/api";
import { useWSHeartbeat } from "../hooks/Heartbeat";
import { log } from "../utils";
import { useInboxSocketHandler } from "../hooks/InboxSocket";
import { useMessageSocketHandler } from "../hooks/MessageSocket";

const SocketContext = createContext<WebSocket | null>(null);

export const useWebSocket = () => useContext(SocketContext);

export const SocketProvider = ({ children }: { children: ReactNode }) => {
    const currentUser = useCurrentUser();
    const [socket, setSocket] = useState<WebSocket | null>(null);

    const onDisconnect = () => {
        if (!socket) return;
        socket.close();
        log("warn", "disconnecting: no heartbeat ack");
    };

    useWSHeartbeat(socket, onDisconnect);
    useInboxSocketHandler(socket);
    useMessageSocketHandler(socket);

    useEffect(() => {
        if (currentUser.isSuccess && !socket) {
            const accessToken = getAccessToken();
            const socket = new WebSocket(
                import.meta.env.VITE_BASE_URL + `/ws?token=${accessToken}`,
            );
            log("info", "connected to WebSocket endpoint");
            setSocket(socket);
        }
    }, [socket, currentUser.isSuccess]);

    useEffect(() => {
        if (!socket) return;

        socket.addEventListener("close", () => {
            log("error", "disconnected from WebSocket endpoint");
            setSocket(null);
        });

        socket.addEventListener("message", (e) => {
            log("debug", `socket message: ${e.data}`);
        });
    }, [socket]);

    return (
        <SocketContext.Provider value={socket}>
            {children}
        </SocketContext.Provider>
    );
};
