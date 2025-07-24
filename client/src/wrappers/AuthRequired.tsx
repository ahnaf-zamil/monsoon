import { useEffect } from "react";
import { useCurrentUser } from "../context/AuthContext";
import { Navigate, Outlet } from "react-router-dom";
import { decryptSeedAndDeriveKeys } from "../crypto/store";
import { CryptoHelper } from "../crypto/helper";
import { useCryptoStore } from "../store/crypto";
import { decodeBase64 } from "tweetnacl-util";

export const AuthRequired = () => {
    const cryptoStore = useCryptoStore();
    const currentUser = useCurrentUser();

    useEffect(() => {
        if (currentUser != null) {
            (async () => {
                console.log("Decrypting keys");
                const seed = await decryptSeedAndDeriveKeys();
                if (seed != null) {
                    const keyPair = CryptoHelper.generateClientKeyPair(seed);
                    cryptoStore.setKeys(
                        decodeBase64(keyPair.ed.priv), // signing key
                        decodeBase64(keyPair.x.priv), // encryption key
                    );
                }
            })();
        }
    }, [currentUser]);

    if (!currentUser) return <Navigate to="/login" />;
    else {
        if (!cryptoStore.hasKeys()) {
            return <></>;
        } else {
            return <Outlet />;
        }
    }
};
