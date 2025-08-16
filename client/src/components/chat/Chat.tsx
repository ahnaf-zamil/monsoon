import type React from "react";
import { MessageBox } from "./MessageBox";
import { sendMessageToConversation } from "../../api/message";
import { useInboxStore } from "../../store/inbox";
import { useCurrentUser } from "../../context/AuthContext";
import { useMessagesForConversation } from "../../hooks/MessageConversation";
import { useRef } from "react";
import { MessageBubble } from "./MessageBubble";
import { Label } from "../ui/label";

export const Chat: React.FC = () => {
  const scrollRef = useRef<HTMLDivElement | null>(null);

  const inboxStore = useInboxStore();
  const currentUser = useCurrentUser();
  const selectedConversation = inboxStore.getSelectedConversation();

  const { data, handleScroll } = useMessagesForConversation(
    selectedConversation,
    scrollRef
  );

  const handleMessageSubmit = async (content: string) => {
    if (!selectedConversation) return;

    const resp = await sendMessageToConversation(
      selectedConversation.conversation_id,
      content
    );
    if (resp.error) {
      console.error(resp.message);
    }
  };

  //h-[calc(100svh)]
  return (
    <>
      <div className="flex-grow sm:block relative overflow-y-hidden">
        <div className="flex items-center h-14 px-5 border-b">
          {/* <Button className="hover:bg-neutral-200 h-11 aspect-square flex items-center justify-center rounded-full p-2.5 sm:hidden">
            <BiArrowBack size={"100%"} />
          </Button> */}
          <h1 className="flex items-center gap-2">
            <>
              <Label className="text-lg">
                {selectedConversation ? selectedConversation.name : "Home"}
              </Label>
            </>
          </h1>
        </div>
        <div className="bg-chatbox fixed top-14 bottom-20 min-h-0 sm:w-[calc(100svw-24rem)] flex flex-col justify-end">
          <div
            className="overflow-y-auto flex-1 px-4 flex py-4  gap-2 flex-col-reverse"
            ref={scrollRef}
            onScroll={handleScroll}
            id="obs"
          >
            {selectedConversation?.conversation_id &&
              data?.map((msg) => {
                const isMyMsg = msg.author_id == currentUser?.data?.id;
                return <MessageBubble msg={msg} isMyMsg={isMyMsg} />;
              })}
          </div>
        </div>

        <MessageBox submitHandler={handleMessageSubmit} />
      </div>
    </>
  );
};
