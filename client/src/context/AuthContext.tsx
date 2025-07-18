import { createContext, useEffect, useState, type ReactNode } from "react";
import type { IUser } from "../types";
import { getAPIAccessToken, getAuthenticatedUser } from "../api/auth";
import { getAccessToken, setAPIAccessToken } from "../api/api";

export const AuthContext = createContext<IUser | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [currentUser, setCurrentUser] = useState<IUser | null>(null);

  useEffect(() => {
    (async () => {
      let accessToken = getAccessToken();
      if (!currentUser && !accessToken) {
        // Only execute if token doesn't exist in state
        let resp = await getAPIAccessToken();
        if (!resp.error) {
          accessToken = resp.data;
          setAPIAccessToken(accessToken);

          resp = await getAuthenticatedUser();
          if (!resp.error) {
            setCurrentUser(resp.data);
          }
        }
      }
    })();
  }, []);

  return (
    <AuthContext.Provider value={currentUser}>{children}</AuthContext.Provider>
  );
};
