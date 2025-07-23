import {
    createContext,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from "react";
import { AuthContext } from "./AuthContext";
import { getAccessToken } from "../api/api";
import { useWSHeartbeat } from "../hooks/Heartbeat";
import { log } from "../utils";

export const SocketContext = createContext<WebSocket | null>(null);

export const useWebSocket = () => useContext(SocketContext);

export const SocketProvider = ({ children }: { children: ReactNode }) => {
    const currentUser = useContext(AuthContext);
    const [socket, setSocket] = useState<WebSocket | null>(null);

    const onDisconnect = () => {
        if (!socket) return;
        socket.close();
        log("warn", "disconnecting: no heartbeat ack");
    };

    useWSHeartbeat(socket, onDisconnect);

    useEffect(() => {
        if (currentUser && !socket) {
            const accessToken = getAccessToken();
            const socket = new WebSocket(
                import.meta.env.VITE_BASE_URL + `/ws?token=${accessToken}`,
            );
            setSocket(socket);
        }
    }, [currentUser]);

    useEffect(() => {
        if (!socket) return;

        socket.addEventListener("close", (_) => {
            log("error", "disconnected from WebSocket endpoint");
            setSocket(null);
        });
    }, [socket]);

    return (
        <SocketContext.Provider value={socket}>
            {children}
        </SocketContext.Provider>
    );
};
