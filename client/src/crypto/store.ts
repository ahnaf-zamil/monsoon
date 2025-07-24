import nacl from "tweetnacl";
import { CryptoHelper } from "./helper";
import localforage from "localforage";

export const storeSeedSecure = async (decryptedSeed: Uint8Array) => {
    /* Encrypt seed using random key generated on client and store in IndexedDB/WebSQL */
    const sessionEncKey = nacl.randomBytes(32);
    const nonce = nacl.randomBytes(12);
    const encryptedSeedLocal = await CryptoHelper.AESGCMEncrypt(
        sessionEncKey,
        nonce,
        decryptedSeed,
    );

    /* Localforage uses IndexedDB or WebSQL on devices which support it, or stores it in Localstorage */
    await localforage.setItem("se", encryptedSeedLocal);
    await localforage.setItem("sn", nonce);

    // Wanted to use Credential API but not all browsers support it, so storing key in IndexedDB. Will need to find better solution, for now storing it here

    // TODO: Reduce attack vectors and vulnerabilities for storing in IndexedDB
    await localforage.setItem("sk", sessionEncKey);
};

export const decryptSeedAndDeriveKeys =
    async (): Promise<Uint8Array | null> => {
        const cipher = await localforage.getItem("se");
        const nonce = await localforage.getItem("sn");
        const key = await localforage.getItem("sk");

        if (key && cipher && nonce) {
            const seed = await CryptoHelper.AESGCMDecrypt(
                key as Uint8Array,
                nonce as Uint8Array,
                cipher as Uint8Array,
            );
            return seed;
        } else {
            // Handle logout and re-authentication
            return null;
        }
    };
