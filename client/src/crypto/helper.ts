import nacl from "tweetnacl";
import { ed25519SeedToX25519PrivateKey } from "./keys";
import { encodeBase64 } from "tweetnacl-util";
import argon2 from "argon2-browser";

export class CryptoHelper {
    static AESGCMEncrypt = async (
        key: Uint8Array,
        nonce: Uint8Array,
        plaintext: Uint8Array,
    ): Promise<Uint8Array> => {
        const cryptoKey = await crypto.subtle.importKey(
            "raw",
            key as BufferSource,
            { name: "AES-GCM" },
            false,
            ["encrypt"],
        );

        const ciphertext = await crypto.subtle.encrypt(
            {
                name: "AES-GCM",
                iv: nonce as BufferSource,
                tagLength: 128, // 16 bytes
            },
            cryptoKey,
            plaintext as BufferSource,
        );

        return new Uint8Array(ciphertext);
    };

    static AESGCMDecrypt = async (
        key: Uint8Array,
        nonce: Uint8Array,
        ciphertext: Uint8Array,
    ): Promise<Uint8Array> => {
        const cryptoKey = await crypto.subtle.importKey(
            "raw",
            key as BufferSource,
            { name: "AES-GCM" },
            false,
            ["decrypt"],
        );

        const decryptedBuffer = await crypto.subtle.decrypt(
            {
                name: "AES-GCM",
                iv: nonce as BufferSource,
                tagLength: 128,
            },
            cryptoKey,
            ciphertext as BufferSource,
        );

        return new Uint8Array(decryptedBuffer);
    };

    static deriveKey = async (
        password: string,
        salt: Uint8Array,
    ): Promise<Uint8Array> => {
        /* Derive Argon2 key using salt and password */
        const res = await argon2.hash({
            pass: password,
            salt: salt,
            time: 1,
            mem: 64 * 1024, // 62 MiB
            hashLen: 32, // bytes
            parallelism: 1,
            type: argon2.ArgonType.Argon2id,
        });
        return res.hash;
    };

    static generateClientKeyPair = (seed: Uint8Array<ArrayBufferLike>) => {
        const edKeyPair = nacl.sign.keyPair.fromSeed(seed);
        const xPriv = ed25519SeedToX25519PrivateKey(seed);
        const xPub = nacl.scalarMult.base(xPriv);

        return {
            ed: {
                pub: encodeBase64(edKeyPair.publicKey),
                priv: encodeBase64(edKeyPair.secretKey),
            },
            x: {
                pub: encodeBase64(xPub),
                priv: encodeBase64(xPriv),
            },
        };
    };
}
