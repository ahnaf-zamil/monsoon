import type React from "react";
import { useInboxStore } from "../store/inbox";
import { InboxEntry } from "./InboxEntry";
import { useCurrentUser } from "../context/AuthContext";
import {
    ConversationSearchbar,
    ConversationSearchComponent,
} from "./ConversationSearch";

import MonsoonLogo from "../static/img/monsoon_logo.png";
import { SidebarCurrentUser } from "./SidebarCurrentUser";
import { useUIStore } from "../store/ui";

const InboxConversations: React.FC = () => {
    const inboxState = useInboxStore();
    return (
        <div className="grid gap-2 overflow-y-auto">
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
        <div className="flex">
            <div className="h-[calc(100svh)] w-full block sm:w-96 sm:block relative border-0 sm:border-r-[1px] sm:border-r-neutral-200 dark:sm:border-r-neutral-800">
                <div className="flex items-center justify-center top-6  w-full absolute">
                    <img
                        src={MonsoonLogo}
                        className="block h-7"
                        draggable="false"
                        alt=""
                    />
                </div>
                <ConversationSearchbar />
                <div className="absolute top-32 left-0 right-0 bottom-0 p-2 flex flex-col justify-between">
                    {uiStore.isSearchingConvo ? (
                        <ConversationSearchComponent />
                    ) : (
                        <InboxConversations />
                    )}
                    <SidebarCurrentUser currentUser={currentUser.data} />
                </div>
            </div>
        </div>
    );
};
