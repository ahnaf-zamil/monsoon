import { create } from "zustand";

interface UIState {
    isSearchingConvo: boolean;
    currentConvoSearchInput?: string;

    setIsSearchingConvo: (v: boolean) => void;
    setCurrentConvoSearchInput: (v?: string) => void;
}

export const useUIStore = create<UIState>((set) => ({
    isSearchingConvo: false,
    setIsSearchingConvo: (v: boolean) =>
        set((state) => ({ ...state, isSearchingConvo: v })),

    setCurrentConvoSearchInput: (v?: string) => {
        set((state) => ({ ...state, currentConvoSearchInput: v }));
    },
}));
