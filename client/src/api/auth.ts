import { apiClient } from "./api";
import type { IAPIResponse, ILoginData } from "./types";

export const fetchUserSalt = async (
    email: string,
): Promise<IAPIResponse<any>> => {
    const response = await apiClient.post("/auth/salt", {
        email,
    });
    return response.data;
};

export const loginUser = async (
    email: string,
    passwordHash: string,
): Promise<IAPIResponse<ILoginData>> => {
    const response = await apiClient.post("/auth/login", {
        email,
        pw_hash: passwordHash,
    });
    return response.data;
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
): Promise<IAPIResponse<any>> => {
    const response = await apiClient.post("/auth/create", {
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
};

export const getAuthenticatedUser = async (): Promise<IAPIResponse<any>> => {
    const response = await apiClient.get("/user/me");
    return response.data;
};

export const getAPIAccessToken = async (): Promise<IAPIResponse<any>> => {
    const response = await apiClient.post("/auth/token");
    return response.data;
};
