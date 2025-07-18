import axios from "axios";

let accessToken: string;

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL + "/api",
  withCredentials: true,
  timeout: 5000,
  headers: { "Content-Type": "application/json" },
  validateStatus: (_) => true // Axios shouldn't handle status codes
});

export const setAPIAccessToken = (token: string) => {
  accessToken = token;
  apiClient.defaults.headers.common["Authorization"] = `Bearer ${accessToken}`;
}

export const getAccessToken = () => {
  return accessToken;
}