import type React from "react";
import { useInboxStore } from "@/store/inbox";
import { InboxEntry } from "../InboxEntry";
import { useCurrentUser } from "@/context/AuthContext";
import {
  ConversationSearchbar,
  ConversationSearchComponent,
} from "./ConversationSearch";

import MonsoonLogo from "@/static/img/monsoon_logo.png";
import { SidebarCurrentUser } from "./SidebarCurrentUser";
import { useUIStore } from "@/store/ui";
import { Separator } from "@/components/ui/separator";

const InboxConversations: React.FC = () => {
  const inboxState = useInboxStore();
  return (
    <div className="flex flex-col h-full gap-4">
      {inboxState.isSynced ? (
        inboxState.conversations.map((convo) => (
          <>
            <InboxEntry
              conversationID={convo.conversation_id}
              name={convo.name}
              last_msg_time={convo.updated_at}
              user_id={convo.user_id}
              key={convo.conversation_id}
            />
          </>
        ))
      ) : (
        <span className="animate-pulse">
          <InboxEntry
            conversationID={""}
            name={""}
            last_msg_time={0}
            user_id={""}
          />
        </span>
      )}
    </div>
  );
};

export const Sidebar: React.FC = () => {
  const currentUser = useCurrentUser();
  const uiStore = useUIStore();

  return (
    <div className="border-r flex flex-col min-h-screen sm:w-96 overflow-y-hidden">
      <div className="flex flex-col flex-1 overflow-y-auto">
        <div className="flex justify-center my-4">
          <img src={MonsoonLogo} className="h-7" draggable="false" alt="" />
        </div>
        <ConversationSearchbar />
        <div className="w-full px-3 flex flex-col">
          <Separator className="my-5" />
          {uiStore.isSearchingConvo ? (
            <ConversationSearchComponent />
          ) : (
            <InboxConversations />
          )}
        </div>
      </div>

      <SidebarCurrentUser currentUser={currentUser.data} />
    </div>
  );
};
