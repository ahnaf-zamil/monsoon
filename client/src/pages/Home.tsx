import type React from "react";
import { Sidebar } from "../components/Sidebar";
import { Chat } from "../components/Chat";
import { useInboxSocketHandler } from "../hooks/InboxSocket";

export const Home: React.FC = () => {
    useInboxSocketHandler();

    return (
        <>
            <div className="flex bg-darkbg max-h-svh">
                <Sidebar />
                <Chat />
            </div>
        </>
    );
};
