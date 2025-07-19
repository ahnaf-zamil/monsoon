import type React from "react";
import MonsoonLogo from "../static/img/monsoon_logo.png";
import { useContext, useEffect, useState } from "react";
import { createUser } from "../api/auth";
import { PiPasswordBold } from "react-icons/pi";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import { FiMail, FiUser } from "react-icons/fi";
import nacl from "tweetnacl";
import { CryptoHelper } from "../crypto/helper";
import { encodeBase64 } from "tweetnacl-util";
import { FaRegAddressCard } from "react-icons/fa";

export const Register: React.FC = () => {
    const navigate = useNavigate();
    const currentUser = useContext(AuthContext);

    const [username, setUsername] = useState<string>("123");
    const [displayName, setDisplayName] = useState<string>("123");
    const [email, setEmail] = useState<string>("ok@ok.ok");
    const [password, setPassword] = useState<string>("123");
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (currentUser) {
            navigate("/");
        }
    }, [currentUser]);

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        try {
            const seed = nacl.randomBytes(32);
            const keys = CryptoHelper.generateClientKeyPair(seed);

            const authSalt = nacl.randomBytes(32);
            const encSalt = nacl.randomBytes(32);

            const pwHash = await CryptoHelper.deriveKey(password, authSalt);
            const encKey = await CryptoHelper.deriveKey(password, encSalt);

            const nonce = nacl.randomBytes(12); // 12 byte nonce
            const encryptedSeed = await CryptoHelper.AESGCMEncrypt(
                encKey,
                nonce,
                seed
            );
            const resp = await createUser(
                username,
                displayName,
                email,
                {
                    sig: keys.ed.pub,
                    enc: keys.x.pub,
                },
                {
                    encSalt: encodeBase64(encSalt),
                    pwSalt: encodeBase64(authSalt),
                },
                encodeBase64(pwHash),
                encodeBase64(encryptedSeed),
                encodeBase64(nonce)
            );
            if (!resp.error) {
                window.location.href = "/";
            } else {
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
            }
        } catch (e) {
            console.error(e);
            setError("An error occured while registering");
        }
    };

    return (
        <div className="flex justify-center items-center flex-col h-[calc(100svh)] dark:bg-black">
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
                            <label htmlFor="username" className="sr-only">
                                Username
                            </label>
                            <div className="flex w-full rounded-md dark:bg-neutral-800 items-center justify-center">
                                <FiUser className="text-xl mx-2 text-primary-text" />
                                <input
                                    type="text"
                                    className="outline-none rounded-md flex-grow pr-4 py-2 placeholder:text-neutral-600 dark:placeholder:text-neutral-500 dark:text-white"
                                    placeholder="Username"
                                    id="username"
                                    required={true}
                                    value={username}
                                    onChange={(e) =>
                                        setUsername(e.target.value)
                                    }
                                />
                            </div>
                        </div>
                        <div className="grid gap-1">
                            <label htmlFor="username" className="sr-only">
                                Display Name
                            </label>
                            <div className="flex w-full rounded-md dark:bg-neutral-800 items-center justify-center">
                                <FaRegAddressCard className="text-xl mx-2 text-primary-text" />
                                <input
                                    type="text"
                                    className="outline-none rounded-md flex-grow pr-4 py-2 placeholder:text-neutral-600 dark:placeholder:text-neutral-500 dark:text-white"
                                    placeholder="Display Name"
                                    id="displayname"
                                    required={true}
                                    value={displayName}
                                    onChange={(e) =>
                                        setDisplayName(e.target.value)
                                    }
                                />
                            </div>
                        </div>
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
                        <button
                            onClick={(e) => handleSubmit(e as any)}
                            className="hover:cursor-pointer hover:bg-primary-darker bg-primary p-2 rounded-md  flex justify-center text-black"
                        >
                            Create Account
                        </button>
                        <p className="text-primary-text text-center">
                            Already have an account?{" "}
                            <a
                                className="cursor-pointer text-secondary-text"
                                onClick={() => navigate("/login")}
                            >
                                Log In
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
