import type React from "react";
import MonsoonLogo from "../static/img/monsoon_logo.png";
import { useContext, useEffect, useState } from "react";
import { fetchUserSalt, loginUser } from "../api/auth";
import { PiPasswordBold } from "react-icons/pi";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import { FiMail } from "react-icons/fi";
import type { IAPIResponse } from "../api/types";
import { decodeBase64, encodeBase64 } from "tweetnacl-util";
import { CryptoHelper } from "../crypto/helper";

export const Login: React.FC = () => {
    const navigate = useNavigate();
    const currentUser = useContext(AuthContext);

    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (currentUser) {
            navigate("/");
        }
    }, [currentUser]);

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        const handleServerError = (resp: IAPIResponse<any>) => {
            switch (resp.status) {
                case 400:
                    setError("Bad request");
                    break;
                case 401:
                    setError("Invalid username/password");
                    break;
                case 500:
                    setError("Internal server error");
                    break;
            }
        };

        try {
            // Fetch user password salt
            let resp = await fetchUserSalt(email);
            if (resp.error) {
                return handleServerError(resp);
            }
            const salt = decodeBase64(resp.data);
            // Derive password hash using salt and password
            const pwHash = await CryptoHelper.deriveKey(password, salt);

            // Send password hash over to server for matching and logging in
            resp = await loginUser(email, encodeBase64(pwHash));
            if (resp.error) {
                return handleServerError(resp);
            }

            let decryptedSeed: Uint8Array;
            try {
                // Derive encryption key from salt and password to use for AES-GCM decryption
                const encKey = await CryptoHelper.deriveKey(
                    password,
                    decodeBase64(resp.data.enc_salt),
                );
                decryptedSeed = await CryptoHelper.AESGCMDecrypt(
                    encKey,
                    decodeBase64(resp.data.nonce),
                    decodeBase64(resp.data.enc_seed),
                );
            } catch (e) {
                return setError("Error while decrypting key seed");
            }
            const keys = CryptoHelper.generateClientKeyPair(decryptedSeed);
            console.log(keys);

            // TODO: Persist keys in session storage and implement inter-tab communication using BroadcastChannel

            window.location.href = "/";
        } catch (e) {
            console.error(e);
            setError("An error occured while logging you in");
        }
    };

    return (
        <div className="flex justify-center items-center flex-col h-[calc(100svh)] dark:bg-black">
            {/* <h1 className="text-teal-400 font-bold text-4xl mb-6">Monsoon</h1> */}
            <img src={MonsoonLogo} className="w-84 pointer-none mb-2" alt="" />
            <div className="flex justify-center items-center ">
                <form
                    className="w-96 p-6 rounded-lg grid gap-2"
                    onSubmit={handleSubmit}
                >
                    {error && (
                        <p className="bg-red-100 border border-red-600 w-full py-2 mb-6 text-center text-red-600 m-auto px-2">
                            {error}
                        </p>
                    )}
                    <div className="grid gap-5">
                        <div className="grid gap-1">
                            <label htmlFor="email" className="sr-only">
                                Email
                            </label>
                            <div className="flex w-full rounded-md dark:bg-neutral-800 items-center justify-center">
                                <FiMail className="text-xl mx-2 text-primary-text" />
                                <input
                                    type="email"
                                    className="outline-none rounded-md flex-grow pr-4 py-2 placeholder:text-neutral-600 dark:placeholder:text-neutral-500 dark:text-white"
                                    placeholder="Email"
                                    id="email"
                                    required={true}
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </div>
                        </div>
                        <div className="grid gap-1">
                            <label htmlFor="password" className="sr-only">
                                Password
                            </label>
                            <div className="flex w-full rounded-md dark:bg-neutral-800 items-center justify-center">
                                <PiPasswordBold className="text-xl mx-2 text-primary-text" />
                                <input
                                    type="password"
                                    className="outline-none rounded-md flex-grow pr-4 py-2 placeholder:text-neutral-600 dark:placeholder:text-neutral-500 dark:text-white"
                                    placeholder="Password"
                                    id="password"
                                    required={true}
                                    value={password}
                                    onChange={(e) =>
                                        setPassword(e.target.value)
                                    }
                                />
                            </div>
                        </div>
                        <button className="hover:cursor-pointer hover:bg-primary-darker bg-primary p-2 rounded-md  flex justify-center text-black">
                            Log in
                        </button>
                        <p className="text-primary-text text-center">
                            Don't have an account?{" "}
                            <a
                                className="cursor-pointer text-secondary-text"
                                onClick={() => navigate("/register")}
                            >
                                Create one
                            </a>
                        </p>
                    </div>
                </form>
            </div>
            <p className="text-center text-gray-400 fixed bottom-[30px] ">
                © K.M Ahnaf Zamil 2025
            </p>
        </div>
    );
};
