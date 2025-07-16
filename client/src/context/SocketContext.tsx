import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { io, Socket } from "socket.io-client";
import { AuthContext } from "./AuthContext";

export const SocketContext = createContext<Socket | undefined>(undefined);

const BASE_URL = import.meta.env.VITE_BASE_URL;

export const SocketProvider = ({ children }: { children: ReactNode }) => {
  const currentUser = useContext(AuthContext);

  const [socket, setSocket] = useState<Socket>();
  const pathnameRef = useRef<string>(location.pathname);

  useEffect(() => {
    pathnameRef.current = location.pathname;
  }, [location]);

  useEffect(() => {
    if (currentUser && !socket) {
      const newSocket = io(BASE_URL, {
        auth: { token: Math.random().toString(16).substring(6) }, // Random 10 ID for now until Auth is implemented
        autoConnect: false,
      });
      setSocket(newSocket);

      return () => {
        newSocket.close();
      };
    }
  }, [currentUser]);

  useEffect(() => {
    if (socket) {
      if (!socket.connected) {
        socket.connect();
      }

      socket.on("connect", () => {
        console.log("Connected to " + BASE_URL + " with ID " + socket.id);
      });
    }
  }, [socket]);

  return (
    <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
  );
};
