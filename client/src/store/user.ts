import { create } from "zustand";
import type { IUser } from "../types";

interface UserCacheState {
    users: {
        [userID: string]: IUser;
    };
    setUser: (userData: IUser) => void;
    getUser: (userID: string) => IUser | undefined;
}

export const useUserCacheStore = create<UserCacheState>((set, get) => ({
    users: {},
    setUser: (userData: IUser) => {
        set((state) => {
            return { users: { ...state.users, [userData.id]: userData } };
        });
    },
    getUser: (userID: string): IUser | undefined => {
        const state = get();
        return state.users[userID];
    },
}));
