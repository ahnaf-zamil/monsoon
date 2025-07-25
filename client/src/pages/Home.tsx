import type React from "react";
import { Sidebar } from "../components/Sidebar";
import { Chat } from "../components/Chat";
import { useInboxStore } from "../store/inbox";
import { useLocation, useParams } from "react-router-dom";
import { useEffect } from "react";

export const Home: React.FC = () => {
    const params = useParams();
    const loc = useLocation()
    const inboxStore = useInboxStore();
    const selectedConversation = inboxStore.getSelectedConversation();

    useEffect(() => {
        if (params.conversationID) {
            inboxStore.setSelectedConvoID(params.conversationID);
        }
    }, [loc]);

    return (
        <>
            <div className="flex bg-darkbg max-h-svh">
                <Sidebar />
                {selectedConversation && <Chat />}
            </div>
        </>
    );
};
