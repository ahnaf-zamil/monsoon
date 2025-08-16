import type React from "react";
import type { IMessageData } from "@/api/types";
import moment from "moment";
import { Card } from "../ui/card";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface Props {
  msg: IMessageData;
  isMyMsg: boolean;
}

const SelectableMessage = ({
  children,
  isOwned,
}: {
  children: React.ReactNode;
  isOwned?: boolean;
}) => {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>{children}</DropdownMenuTrigger>
      <DropdownMenuContent className="ring-0 outline-none" align="start">
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Add Reaction</DropdownMenuSubTrigger>
            <DropdownMenuSubContent>
              <DropdownMenuItem>👍</DropdownMenuItem>
              <DropdownMenuItem>😂</DropdownMenuItem>
              <DropdownMenuItem>😭</DropdownMenuItem>
            </DropdownMenuSubContent>
          </DropdownMenuSub>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem>Reply</DropdownMenuItem>
        <DropdownMenuItem>Copy Text</DropdownMenuItem>
        <DropdownMenuSeparator />
        {!isOwned && (
          <DropdownMenuItem className="text-red-500">
            Report Message
          </DropdownMenuItem>
        )}
        {isOwned && (
          <DropdownMenuItem className="text-red-500">
            Delete Message
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export const MessageBubble: React.FC<Props> = ({ msg, isMyMsg }) => {
  const msgMoment = moment.unix(msg.created_at);
  const today = moment();
  const isToday = msgMoment.isSame(today, "day");

  return (
    <Card className="py-3 border-none">
      <small
        className={`m-3 text-foreground/50 z-0 flex justify-${
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
          <SelectableMessage isOwned>
            <div className={`bg-foreground/20 rounded-2xl px-4 py-2 max-w-xs`}>
              {msg.content}
            </div>
          </SelectableMessage>
        ) : (
          <SelectableMessage>
            <div
              className={`bg-accent px-4 py-2 rounded-2xl max-w-xs text-start`}
            >
              {msg.content}
            </div>
          </SelectableMessage>
        )}
      </div>
    </Card>
  );
};
