import type React from "react";
import { Sidebar } from "../components/Sidebar";
import { Chat } from "../components/Chat";
import { useInboxStore } from "../store/inbox";

export const Home: React.FC = () => {
    const inboxStore = useInboxStore();
    const selectedConversation = inboxStore.getSelectedConversation();

    return (
        <>
            <div className="flex bg-darkbg max-h-svh">
                <Sidebar />
                {selectedConversation && <Chat />}
            </div>
        </>
    );
};
