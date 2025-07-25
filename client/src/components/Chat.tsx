import type React from "react";
import { BiArrowBack } from "react-icons/bi";
import { MessageBox } from "./MessageBox";
import {
    fetchConversationMessages,
    sendMessageToConversation,
} from "../api/message";
import { useInboxStore } from "../store/inbox";
import { useEffect } from "react";
import { useMessageStore } from "../store/message";
import { useCurrentUser } from "../context/AuthContext";
import { log } from "../utils";
import { useUserCacheStore } from "../store/user";
import { useQuery } from "@tanstack/react-query";
import type { IMessageData } from "../api/types";

export const Chat: React.FC = () => {
    const inboxStore = useInboxStore();
    const userCache = useUserCacheStore();
    const messageStore = useMessageStore();
    const currentUser = useCurrentUser();
    const selectedConversation = inboxStore.getSelectedConversation();

    const messageQuery = useQuery({
        queryKey: ["chat-messages"],
        queryFn: async (): Promise<IMessageData[]> => {
            const resp = await fetchConversationMessages(
                selectedConversation!.conversation_id,
                20,
            );
            if (resp.error) {
                throw Error(resp.message);
            }
            return resp.data;
        },
        retry: true,
        enabled: false,
    });

    const handleMessageSubmit = async (content: string) => {
        if (!selectedConversation) return;

        const resp = await sendMessageToConversation(
            selectedConversation.conversation_id,
            content,
        );
        if (resp.error) {
            console.error(resp.message);
        }
    };

    const fetchAndCacheConversationUser = async (
        conversationID: string,
        userID: string,
    ) => {
        // WIP: Fetch DM user data from API and store in cache
        const cachedUser = userCache.getUser(userID);

        if (!cachedUser) {
            console.log("User data not cached");
        } else {
            console.log("User data is cached");
        }
    };

    useEffect(() => {
        if (selectedConversation != null) {
            const convoID = selectedConversation.conversation_id;

            if (selectedConversation.type === "DM") {
                fetchAndCacheConversationUser(
                    selectedConversation.conversation_id,
                    selectedConversation.user_id!,
                ); // user_id will always be present in DM conversations
            }

            const msg = messageStore.getConversationMessages(convoID);
            if (msg == undefined) {
                log("debug", `convo ${convoID} not cached, fetching messages`);

                messageQuery.refetch();
            } else {
                log("debug", `convo ${convoID} has cached messages`);
            }
        }
    }, [selectedConversation]);

    useEffect(() => {
        if (selectedConversation && messageQuery.isSuccess) {
            messageStore.storeMessages(
                selectedConversation?.conversation_id,
                messageQuery.data!,
            );
        }
    }, [messageQuery.isSuccess, messageQuery.isRefetching]);

    return (
        <>
            <div className="flex-grow sm:block relative h-[calc(100svh)]">
                <div className="flex items-center gap-3 absolute top-0 right-0 left-0 h-14 px-5 sm:px-10 border border-b-neutral-200 border-x-0 border-t-0 dark:border-b-neutral-800">
                    <button className="hover:bg-neutral-200 h-11 aspect-square flex items-center justify-center rounded-full p-2.5 sm:hidden dark:text-white dark:hover:bg-neutral-800">
                        <BiArrowBack size={"100%"} />
                    </button>
                    <h1 className="text-lg dark:text-white flex items-center gap-2">
                        <>
                            <p>
                                {selectedConversation
                                    ? selectedConversation.name
                                    : "Home"}
                            </p>
                        </>
                    </h1>
                </div>
                {selectedConversation && (
                    <div className="bg-chatbox fixed top-14 bottom-20 min-h-0 sm:w-[calc(100svw-24rem)] flex flex-col justify-end">
                        {/* Scrollable message container with reverse column */}
                        <div className="overflow-y-auto flex-1 px-4 flex   gap-2 flex-col-reverse">
                            {messageStore
                                .getConversationMessages(
                                    selectedConversation.conversation_id,
                                )
                                ?.map((msg) => {
                                    const isMyMsg =
                                        msg.author_id == currentUser?.data?.id;
                                    return (
                                        <div
                                            key={msg.id}
                                            className={`flex justify-${
                                                !isMyMsg ? "start" : "end"
                                            }`}
                                        >
                                            {isMyMsg ? (
                                                <div
                                                    className={`bg-chatbubble-sender text-white px-4 py-2 rounded-lg max-w-xs`}
                                                >
                                                    {msg.content}
                                                </div>
                                            ) : (
                                                <div
                                                    className={`bg-chatbubble-recipient text-white px-4 py-2 rounded-lg max-w-xs`}
                                                >
                                                    {msg.content}
                                                </div>
                                            )}
                                        </div>
                                    );
                                })}
                        </div>
                    </div>
                )}

                <MessageBox submitHandler={handleMessageSubmit} />
            </div>
        </>
    );
};
