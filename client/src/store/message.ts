import { create } from "zustand";
import type { IMessageData } from "../api/types";

interface MessageState {
    messages: {
        [conversationID: string]: IMessageData[];
    };
    storeMessages: (
        conversationID: string,
        newMessages: IMessageData[]
    ) => void;
    getConversationMessages: (conversationID: string) => IMessageData[] | undefined;
}

export const useMessageStore = create<MessageState>((set, get) => ({
    messages: {},
    storeMessages: (conversationID: string, newMessages: IMessageData[]) => {
        set((state) => {
            let msg: IMessageData[];
            const existingMessages = state.messages[conversationID];
            if (existingMessages == undefined) {
                msg = newMessages;
            } else {
                msg = [...existingMessages, ...newMessages];
            }
            msg.sort((a, b) => b.created_at - a.created_at);
            return {
                messages: { ...state.messages, [conversationID]: msg },
            };
        });
    },
    getConversationMessages: (conversationID: string): IMessageData[] | undefined => {
        const state = get();
        const messages = state.messages[conversationID];
        return messages;
    },
}));
