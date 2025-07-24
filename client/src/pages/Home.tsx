import type React from "react";
import { Sidebar } from "../components/Sidebar";
import { Chat } from "../components/Chat";

export const Home: React.FC = () => {
    return (
        <>
            <div className="flex bg-darkbg max-h-svh">
                <Sidebar />
                <Chat />
            </div>
        </>
    );
};
