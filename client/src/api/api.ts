import axios from "axios";

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL + "/api",
  withCredentials: true,
  timeout: 5000,
  headers: { "Content-Type": "application/json" },
});
