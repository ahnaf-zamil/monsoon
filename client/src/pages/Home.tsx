import type React from "react";
import { Sidebar } from "../components/Sidebar";
import { Chat } from "../components/Chat";
import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";

export const Home: React.FC = () => {
  const user = useContext(AuthContext);

  return (
    <>
      <div className="flex bg-darkbg max-h-svh">
        <Sidebar />
        <Chat />
      </div>
    </>
  );
};
