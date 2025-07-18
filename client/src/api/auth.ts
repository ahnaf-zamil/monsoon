import { apiClient } from "./api";
import type { IAPIResponse } from "./types";

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

export const getAuthenticatedUser = async (): Promise<IAPIResponse> => {
  // TODO: Implement on backend
  const response = await apiClient.get("/user/me");
  return response.data;
};

export const getAPIAccessToken = async(): Promise<IAPIResponse> => {
  const response = await apiClient.post("/user/token");
  return response.data;
}