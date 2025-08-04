import type React from "react";
import type { IMessageData } from "../api/types";
import moment from "moment";

interface Props {
    msg: IMessageData;
    isMyMsg: boolean;
}

export const MessageBubble: React.FC<Props> = ({ msg, isMyMsg }) => {
    const msgMoment = moment.unix(msg.created_at);
    const today = moment();
    const isToday = msgMoment.isSame(today, "day");

    return (
        <>
            <small
                className={`mx-1 text-slate-400 flex justify-${!isMyMsg ? "start" : "end"}`}
            >
                {isToday
                    ? "Today at " + msgMoment.format("h:mm A")
                    : msgMoment.format("MMMM Do YYYY | h:mm A")}
            </small>
            <div
                key={msg.id}
                className={`flex justify-${!isMyMsg ? "start" : "end"}`}
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
        </>
    );
};
