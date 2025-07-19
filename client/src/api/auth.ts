import { apiClient } from "./api";
import type { IAPIResponse } from "./types";

export const loginUser = async (
    email: string,
    password: string,
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

export const createUser = async (
    username: string,
    displayName: string,
    email: string,
    pub_keys: { sig: string; enc: string },
    salts: {
        pwSalt: string;
        encSalt: string;
    },
    pwHash: string,
    encSeed: string,
    nonce: string,
): Promise<any> => {
    try {
        const response = await apiClient.post("/user/create", {
            username,
            display_name: displayName,
            email,
            pub_keys: pub_keys,
            salts: {
                enc_salt: salts.encSalt,
                pw_salt: salts.pwSalt,
            },
            pw_hash: pwHash,
            enc_seed: encSeed,
            nonce: nonce,
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

export const getAPIAccessToken = async (): Promise<IAPIResponse> => {
    const response = await apiClient.post("/user/token");
    return response.data;
};
