import type { IUser } from "../types";
import { apiClient } from "./api";

export const loginUser = async (username: string, password: string): Promise<any> => {
  try {
    const response = await apiClient.post("/auth/login", {
      username,
      password,
    });
    return response.data;
  } catch (error) {
    return error;
  }
};


export const getAuthenticatedUser = async (): Promise<IUser | null> => {
    try {
        const response = await apiClient.get("/auth/me");
        return response.data;
    }
    catch (error) {
        return null;
    }
}