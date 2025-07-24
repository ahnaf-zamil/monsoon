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

export const Chat: React.FC = () => {
    const inboxStore = useInboxStore();
    const messageStore = useMessageStore();
    const currentUser = useCurrentUser();
    const selectedConversation = inboxStore.getSelectedConversation();

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

    useEffect(() => {
        if (selectedConversation != null) {
            (async () => {
                const convoID = selectedConversation.conversation_id;

                const msg = messageStore.getConversationMessages(convoID);
                if (msg == undefined) {
                    log(
                        "debug",
                        `convo ${convoID} not cached, fetching messages`,
                    );
                    const resp = await fetchConversationMessages(
                        selectedConversation.conversation_id,
                        20,
                    );

                    if (!resp.error) {
                        messageStore.storeMessages(convoID, resp.data);
                    } else {
                        console.error(resp.error);
                    }
                } else {
                    log("debug", `convo ${convoID} has cached messages`);
                }
            })();
        }
    }, [selectedConversation]);

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
                                ?.map((msg, _) => {
                                    const isMyMsg =
                                        msg.author_id == currentUser?.id;
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
