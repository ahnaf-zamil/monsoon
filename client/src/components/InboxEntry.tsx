import type React from "react";
import { formatUnixToLocalTime } from "../utils";
import { useInboxStore } from "../store/inbox";
import { useNavigate } from "react-router-dom";

interface Props {
    conversationID: string;
    name: string;
    last_msg_time: number;
    user_id: string | null;
}

export const InboxEntry: React.FC<Props> = (props) => {
    const inboxStore = useInboxStore();
    const navigate = useNavigate()
    const selectedConversation = inboxStore.getSelectedConversation();

    return (
        <div
            onClick={() => {
                navigate(`/conversations/${props.conversationID}`)
            }}
            className={
                "hover:bg-neutral-900 rounded-md flex gap-3 p-3 items-center justify-between cursor-pointer " +
                (selectedConversation?.conversation_id == props.conversationID
                    ? "bg-neutral-800"
                    : "")
            }
        >
            <div className="flex gap-3 items-center">
                <div className="relative">
                    <div className="w-12 h-12 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center bg-purple-100 text-purple-700">
                        <img
                            src={undefined}
                            alt=""
                            className="object-cover w-full h-full"
                        />
                    </div>
                    <div className="absolute bg-green-500 rounded-md w-3 h-3 bottom-0 right-0"></div>
                </div>
                <div className="grid items-center">
                    <h2 className="text-lg dark:text-white flex items-center gap-2">
                        {props.name}
                    </h2>
                    <p className="text-neutral-600 text-sm truncate dark:text-neutral-500">
                        Last msg...
                    </p>
                </div>
            </div>
            <div className="flex flex-col-reverse items-end self-end gap-1 text-sm text-neutral-600 flex-shrink-0 dark:text-neutral-500">
                <span>{formatUnixToLocalTime(props.last_msg_time)}</span>
            </div>
        </div>
    );
};
