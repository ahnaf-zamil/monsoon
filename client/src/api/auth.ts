import type { IUser } from "../types";
import { apiClient } from "./api";

export const loginUser = async (
  email: string,
  password: string
): Promise<any> => {
  try {
    const response = await apiClient.post("/user/login", {
      email,
      password,
    });
    return response.data;
  } catch (error) {
    return error;
  }
};

export const getAuthenticatedUser = async (): Promise<IUser | null> => {
  // TODO: Implement on backend
  try {
    const response = await apiClient.get("/user/@me");
    return response.data;
  } catch (error) {
    return null;
  }
};
