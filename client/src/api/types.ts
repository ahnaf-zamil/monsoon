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
