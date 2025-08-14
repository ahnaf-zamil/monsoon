import type React from "react";
import { formatUnixToLocalTime } from "../../utils";
import { useInboxStore } from "../../store/inbox";
import { useNavigate } from "react-router-dom";

import { Card, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

interface Props {
  conversationID: string;
  name: string;
  last_msg_time: number;
  user_id: string | null;
  onClick?: () => void;
}

export const InboxEntry: React.FC<Props> = (props) => {
  const inboxStore = useInboxStore();
  const navigate = useNavigate();
  const selectedConversation = inboxStore.getSelectedConversation();

  return (
    <Card
      onClick={() => {
        if (props.onClick) props.onClick();
        navigate(`/conversations/${props.conversationID}`);
      }}
      className={`rounded-xl hover:cursor-pointer transition ease 1s hover:bg-foreground/20 h-24 flex items-center ${
        selectedConversation?.conversation_id == props.conversationID
          ? "bg-foreground/10"
          : ""
      }`}
    >
      <div className="w-full px-3">
        <div className="flex flex-row gap-4 items-center">
          <Avatar className="w-16 h-16">
            <AvatarImage src="https://github.com/ahnaf-zamil.png" alt="" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
          <div className="flex flex-col w-full">
            <CardTitle>
              <Label className="text-lg font-semibold">{props.name}</Label>
            </CardTitle>
            <div className="w-full">
              <div className="flex flex-row justify-between">
                <Label className="text-foreground/40 font-normal w-28 truncate">
                  The quick brown fox jumped over the poop
                </Label>
                <Label className="text-foreground/30">
                  {formatUnixToLocalTime(props.last_msg_time)}
                </Label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
};
