import { sha512 } from "@noble/hashes/sha2";

function clampPrivateKey(keyBytes: Uint8Array<ArrayBuffer>) {
    keyBytes[0] &= 248;
    keyBytes[31] &= 127;
    keyBytes[31] |= 64;
    return keyBytes;
}

export function ed25519SeedToX25519PrivateKey(
    seed: Uint8Array<ArrayBufferLike>,
) {
    const hashed = sha512(seed);
    const xPriv = new Uint8Array(hashed.slice(0, 32));
    return clampPrivateKey(xPriv);
}
