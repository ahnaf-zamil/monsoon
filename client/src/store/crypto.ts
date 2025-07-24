import { create } from "zustand";

interface CryptoState {
    signingKey: Uint8Array | null;
    encryptionKey: Uint8Array | null;
    setKeys: (signingKey: Uint8Array, encryptionKey: Uint8Array) => void;
    hasKeys: () => boolean;
}

export const useCryptoStore = create<CryptoState>((set, get) => ({
    signingKey: null,
    encryptionKey: null,

    setKeys: (signingKey: Uint8Array, encryptionKey: Uint8Array) => {
        console.log("Storing keys in state", signingKey, encryptionKey);
        set((state) => ({ ...state, signingKey, encryptionKey }));
    },

    hasKeys: () => {
        const state = get();
        if (state.signingKey != null && state.encryptionKey != null) {
            return true;
        }
        return false;
    },
}));
