import { create } from "zustand";
import type { IMessageData } from "../api/types";

interface MessageState {
    messages: {
        [conversationID: string]: IMessageData[];
    };
    storeMessages: (
        conversationID: string,
        newMessages: IMessageData[],
    ) => void;
    getConversationMessages: (
        conversationID: string,
    ) => IMessageData[] | undefined;
}

export const useMessageStore = create<MessageState>((set, get) => ({
    messages: {},
    storeMessages: (conversationID: string, newMessages: IMessageData[]) => {
    set((state) => {
        const existingMessages = state.messages[conversationID] || [];

        // Merge and deduplicate based on message ID
        const allMessages = [...existingMessages, ...newMessages];

        const uniqueMessagesMap: Record<string, IMessageData> = {};
        for (const msg of allMessages) {
            uniqueMessagesMap[msg.id] = msg;
        }

        const dedupedMessages = Object.values(uniqueMessagesMap);
        dedupedMessages.sort((a, b) => b.created_at - a.created_at);

        return {
            messages: {
                ...state.messages,
                [conversationID]: dedupedMessages,
            },
        };
    });
},
    getConversationMessages: (
        conversationID: string,
    ): IMessageData[] | undefined => {
        const state = get();
        const messages = state.messages[conversationID];
        return messages;
    },
}));
