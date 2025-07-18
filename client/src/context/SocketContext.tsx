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

export const SocketContext = createContext<WebSocket | null>(null);

export const SocketProvider = ({ children }: { children: ReactNode }) => {
  const currentUser = useContext(AuthContext);
  const [socket, setSocket] = useState<WebSocket|null>(null);

  useWSHeartbeat(socket);

  useEffect(() => {
    if (currentUser && !socket) {
      const accessToken = getAccessToken();
      const socket = new WebSocket(
        import.meta.env.VITE_BASE_URL + `/ws?token=${accessToken}`
      );
      setSocket(socket);
    }
  }, [currentUser]);
  return (
    <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
  );
};
