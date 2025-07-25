import {
    createContext,
    useContext,
    useEffect,
    useRef,
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

const reconnectInterval = 4000; // 4 seconds

export const SocketProvider = ({ children }: { children: ReactNode }) => {
    const currentUser = useCurrentUser();
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [connected, setConnected] = useState<boolean>(false);
    const hasConnectedOnceRef = useRef<boolean>(false);

    const onNoHeartbeatAck = () => {
        log("warn", "disconnecting: no heartbeat ack");
        socket?.close();
    };

    useWSHeartbeat(socket, connected, onNoHeartbeatAck);
    useInboxSocketHandler(socket);
    useMessageSocketHandler(socket);
    const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(
        null,
    );

    useEffect(() => {
        // clear the reconnect timeout on mount
        return () => {
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
        };
    }, []);

    useEffect(() => {
        if (!currentUser.isSuccess || connected || socket) return;
        const attemptConnection = () => {
            const accessToken = getAccessToken();
            const newSocket = new WebSocket(
                import.meta.env.VITE_BASE_URL + `/ws?token=${accessToken}`,
            );
            log("info", "attempting connection to WebSocket endpoint");

            setSocket(newSocket);
            setConnected(true);

            reconnectTimeoutRef.current = null;
            hasConnectedOnceRef.current = true;
        };

        if (hasConnectedOnceRef.current) {
            // if already connected once then this time attempt delayed reconnect
            if (!reconnectTimeoutRef.current) {
                reconnectTimeoutRef.current = setTimeout(
                    attemptConnection,
                    reconnectInterval,
                ); // delay on reconnect
            }
        } else {
            // connect asap on first time
            attemptConnection();
        }
    }, [connected, currentUser.isSuccess, socket]);

    useEffect(() => {
        const handleOpen = () => {
            log("info", "connected to WebSocket endpoint");
        };

        const handleClose = () => {
            log("error", "disconnected from WebSocket endpoint");
            setSocket(null);
            setConnected(false);
        };

        const handleMessage = (e: MessageEvent) => {
            log("debug", `socket message: ${e.data}`);
        };

        socket?.addEventListener("close", handleClose);
        socket?.addEventListener("open", handleOpen);
        socket?.addEventListener("message", handleMessage);

        return () => {
            socket?.removeEventListener("open", handleOpen);
            socket?.removeEventListener("close", handleClose);
            socket?.removeEventListener("message", handleMessage);
        };
    }, [connected, socket]);

    return (
        <SocketContext.Provider value={socket}>
            {children}
        </SocketContext.Provider>
    );
};
