import { create } from "zustand";
import { log } from "../utils";

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
        log("info", "stored encryption and signing keys in memory");
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
