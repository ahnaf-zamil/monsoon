import { apiClient } from "./api";

export const searchUser = async (username: string) => {
  const response = await apiClient.get(`/user/search/${username}`);

  if (response.status != 200) {
    return null;
  }

  return response.data;
};
