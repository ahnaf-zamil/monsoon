import {
    createContext,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from "react";
import type { IUser } from "../types";
import { getAPIAccessToken, getAuthenticatedUser } from "../api/auth";
import { getAccessToken, setAPIAccessToken } from "../api/api";
import { useQuery } from "@tanstack/react-query";

interface IUserState {
    isPending: boolean;
    isError: boolean;
    isSuccess: boolean;
    error: Error | null;
    data: IUser | undefined;
}

const defaults = {
    isPending: true,
    isError: false,
    isSuccess: false,
    error: null,
    data: undefined,
};

export const AuthContext = createContext<IUserState>(defaults);

export const useCurrentUser = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [currentUser, setCurrentUser] = useState<IUserState>(defaults);

    const { isPending, isError, error, data, isSuccess, refetch } = useQuery({
        queryKey: ["current-user"],
        queryFn: async (): Promise<IUser | undefined> => {
            let accessToken = getAccessToken();
            if (!accessToken) {
                const resp = await getAPIAccessToken();
                if (resp.error) {
                    throw Error(resp.message);
                }
                accessToken = resp.data;
                setAPIAccessToken(accessToken);
            }
            const resp = await getAuthenticatedUser();
            if (resp.error) {
                throw Error(resp.message);
            }
            return resp.data;
        },
        retry: false,
        enabled: false,
        refetchInterval: false,
    });

    useEffect(() => {
        setCurrentUser({
            isPending,
            isError,
            error,
            data,
            isSuccess,
        });
    }, [isPending, isError, error, data, isSuccess]);

    useEffect(() => {
        refetch();
    }, []);

    return (
        <AuthContext.Provider value={currentUser}>
            {children}
        </AuthContext.Provider>
    );
};
