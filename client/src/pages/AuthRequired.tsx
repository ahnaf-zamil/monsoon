import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import { Navigate, Outlet } from "react-router-dom";

export const AuthRequired = () => {
  const currentUser = useContext(AuthContext);
  return currentUser ? <Outlet /> : <Navigate to="/login" />;
};
