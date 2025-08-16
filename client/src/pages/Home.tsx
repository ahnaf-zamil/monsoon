import type React from "react";
import { Sidebar } from "../components/chat/sidebar/Sidebar";
import { Chat } from "../components/chat/Chat";
import { useInboxStore } from "../store/inbox";
import { useLocation, useParams } from "react-router-dom";
import { useEffect } from "react";

export const Home: React.FC = () => {
  const params = useParams();
  const loc = useLocation();
  const inboxStore = useInboxStore();
  const selectedConversation = inboxStore.getSelectedConversation();

  useEffect(() => {
    if (params.conversationID) {
      inboxStore.setSelectedConvoID(params.conversationID);
    }
  }, [loc]);

  return (
    <>
      <div className="flex min-h-screen">
        <Sidebar />
        {selectedConversation && <Chat />}
      </div>
    </>
  );
};
