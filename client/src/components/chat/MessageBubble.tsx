import type React from "react";
import type { IMessageData } from "@/api/types";
import moment from "moment";
import { Card } from "../ui/card";

interface Props {
  msg: IMessageData;
  isMyMsg: boolean;
}

export const MessageBubble: React.FC<Props> = ({ msg, isMyMsg }) => {
  const msgMoment = moment.unix(msg.created_at);
  const today = moment();
  const isToday = msgMoment.isSame(today, "day");

  return (
    <Card className="py-3 border-none">
      <small
        className={`mx-3 text-foreground/50 flex justify-${
          !isMyMsg ? "start" : "end"
        }`}
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
          <div className={`bg-foreground/20 rounded-2xl px-4 py-2 max-w-xs`}>
            {msg.content}
          </div>
        ) : (
          <div className={`bg-accent px-4 py-2 rounded-2xl max-w-xs`}>
            {msg.content}
          </div>
        )}
      </div>
    </Card>
  );
};
