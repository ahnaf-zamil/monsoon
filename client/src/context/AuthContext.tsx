import { createContext, useEffect, useState, type ReactNode } from "react";
import type { IUser } from "../types";
import { getAuthenticatedUser } from "../api/auth";
import { useLocation } from "react-router-dom";

export const AuthContext = createContext<IUser | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<IUser | null>(null);

  const location = useLocation();

  useEffect(() => {
    (async () => {
      const currentUser = await getAuthenticatedUser();
      setUser(currentUser);
    })();
  }, [location.pathname]);

  return <AuthContext.Provider value={user}>{children}</AuthContext.Provider>;
};
