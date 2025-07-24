import { create } from "zustand";
import type { IInboxEntry } from "../ws/types";

interface InboxState {
    isSynced: boolean;
    conversations: IInboxEntry[];
    syncConversations: (convos: IInboxEntry[]) => void;

    selectedConvoID: string | null;
    setSelectedConvoID: (convoID: string) => void;
    getSelectedConversation: () => IInboxEntry | undefined;
}

export const useInboxStore = create<InboxState>((set, get) => ({
    isSynced: false,
    conversations: [],
    selectedConvoID: null,

    syncConversations: (convos: IInboxEntry[]) => {
        set((state) => ({ ...state, conversations: convos, isSynced: true }));
    },
    setSelectedConvoID: (convoID: string) =>
        set((state) => ({ ...state, selectedConvoID: convoID })),

    getSelectedConversation: () => {
        const state = get();
        if (state.selectedConvoID == null) {
            return;
        }

        for (let i = 0; i < state.conversations.length; i++) {
            const c = state.conversations[i];
            if (c.conversation_id == state.selectedConvoID) {
                return c;
            }
        }
    },
}));
