import { apiClient } from "./api";
import type { IAPIResponse } from "./types";

export const sendDirectMessageToUser = async (
    userID: string, content: string
): Promise<IAPIResponse<any>> => {
    const response = await apiClient.post(`/message/user/${userID}`, {
        content
    });
    return response.data;
};
