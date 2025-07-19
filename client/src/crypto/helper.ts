import nacl from "tweetnacl";
import { ed25519SeedToX25519PrivateKey } from "./keys";
import { encodeBase64 } from "tweetnacl-util";

export class CryptoHelper {
    static importAESKey = async (rawKey: ArrayBuffer): Promise<CryptoKey> => {
        return await crypto.subtle.importKey(
            "raw",
            rawKey,
            { name: "AES-GCM" },
            false,
            ["encrypt", "decrypt"],
        );
    };

    static AESGCMEncrypt = async (
        cryptoKey: CryptoKey,
        nonce: Uint8Array,
        plaintext: Uint8Array,
    ): Promise<Uint8Array> => {
        const encryptedBuffer = await crypto.subtle.encrypt(
            {
                name: "AES-GCM",
                iv: nonce as BufferSource,
                tagLength: 128,
            },
            cryptoKey,
            plaintext as BufferSource,
        );

        return new Uint8Array(encryptedBuffer);
    };

    static deriveKey = async (password: string, salt: Uint8Array) => {
        const keyMaterial = await window.crypto.subtle.importKey(
            "raw",
            new TextEncoder().encode(password),
            { name: "PBKDF2" },
            false,
            ["deriveBits"],
        );

        const derivedBits = await window.crypto.subtle.deriveBits(
            {
                name: "PBKDF2",
                salt: salt as BufferSource,
                iterations: 150000,
                hash: "SHA-256",
            },
            keyMaterial,
            32 * 8,
        );

        return new Uint8Array(derivedBits);
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
