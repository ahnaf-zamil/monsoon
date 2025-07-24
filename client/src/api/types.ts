export interface IAPIResponse<T> {
    error: boolean;
    message: string;
    data: T;
    status: number;
}

export interface ILoginData {
    enc_salt: string;
    enc_seed: string;
    nonce: string;
}

export interface IMessageData {
    id: string;
    conversation_id: string;
    author_id: string;
    content: string | null;
    created_at: number;
    edited_at: number;
}