import type React from "react";
import { useInboxStore } from "../store/inbox";
import { InboxEntry } from "./InboxEntry";
import { useCurrentUser } from "../context/AuthContext";

export const Sidebar: React.FC = () => {
    const inboxState = useInboxStore();

    const currentUser = useCurrentUser();

    return (
        <>
            <div className="h-[calc(100svh)] w-full block sm:w-96 sm:block relative border-0 sm:border-r-[1px] sm:border-r-neutral-200 dark:sm:border-r-neutral-800">
                <div className="absolute top-0 left-0 right-0 bottom-0 p-2 flex flex-col justify-between">
                    <div className="grid gap-2 overflow-y-auto">
                        {inboxState.isSynced ? (
                            inboxState.conversations.map((convo) => (
                                <InboxEntry
                                    conversationID={convo.conversation_id}
                                    name={convo.name}
                                    last_msg_time={convo.updated_at}
                                    user_id={convo.user_id}
                                    key={convo.conversation_id}
                                />
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
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                            <div className="w-12 aspect-square rounded-full overflow-hidden">
                                <img
                                    src={undefined}
                                    alt=""
                                    className="object-cover w-full h-full"
                                />
                            </div>
                            <div className="grid items-center">
                                <>
                                    <p className="text-base leading-4 dark:text-white">
                                        {currentUser?.data?.display_name}
                                    </p>
                                    <p className="text-neutral-600 text-sm leading-4 dark:text-neutral-500">
                                        @{currentUser?.data?.username}
                                    </p>
                                </>
                            </div>
                        </div>
                        <button
                            name="logout"
                            className="hover:bg-neutral-300 h-fit p-3 rounded-full bg-neutral-200 dark:bg-neutral-800 dark:text-white dark:hover:bg-neutral-900"
                        ></button>
                    </div>
                </div>
            </div>
        </>
    );
};
