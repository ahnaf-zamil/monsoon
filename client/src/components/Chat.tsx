import type React from "react";
import { BiArrowBack } from "react-icons/bi";
import { MessageBox } from "./MessageBox";
import { sendMessageToConversation } from "../api/message";
import { useInboxStore } from "../store/inbox";

export const Chat: React.FC = () => {
    const inboxStore = useInboxStore();
    const selectedConversation = inboxStore.getSelectedConversation();

    const handleMessageSubmit = async (content: string) => {
        if (!selectedConversation) return;

        const resp = await sendMessageToConversation(
            selectedConversation.conversation_id,
            content
        );
        if (!resp.error) {
            console.log(resp.data);
        } else {
            console.error(resp.message);
        }
    };

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
                <div className="fixed top-14 bottom-20 min-h-0 w-full flex flex-col justify-end">
                    <div className="grid gap-2 p-2 pb-8 overflow-y-auto relative"></div>
                </div>
                <MessageBox submitHandler={handleMessageSubmit} />
            </div>
        </>
    );
};
