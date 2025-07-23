import { useCurrentUser } from "../context/AuthContext";
import { Navigate, Outlet } from "react-router-dom";

export const AuthRequired = () => {
    const currentUser = useCurrentUser();
    return currentUser ? <Outlet /> : <Navigate to="/login" />;
};
