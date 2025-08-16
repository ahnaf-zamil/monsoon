import MonsoonLogo from "../../static/img/monsoon_logo.png";
import React, { useEffect, useState } from "react";
import { FiMail, FiUser } from "react-icons/fi";
import { PiPasswordBold } from "react-icons/pi";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { toast } from "@/hooks/use-toast";
import { createUser, fetchUserSalt, loginUser } from "@/api/auth";
import type { IAPIResponse } from "@/api/types";
import { useCurrentUser } from "@/context/AuthContext";
import { CryptoHelper } from "../../crypto/helper";
import { storeSeedSecure } from "@/crypto/store";
import { LoadingPage } from "@/pages/Loading";
import { useNavigate, Navigate } from "react-router-dom";
import { decodeBase64, encodeBase64 } from "tweetnacl-util";
import { FaRegAddressCard } from "react-icons/fa";
import nacl from "tweetnacl";

type FormAuthProps = {
  isLogin: boolean;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => Promise<void>;

  setEmail: React.Dispatch<React.SetStateAction<string>>;
  setPassword: React.Dispatch<React.SetStateAction<string>>;

  setUsername?: React.Dispatch<React.SetStateAction<string>>;
  setDisplayName?: React.Dispatch<React.SetStateAction<string>>;

  setError: React.Dispatch<React.SetStateAction<string | null>>;

  email: string;
  password: string;

  username?: string;
  displayName?: string;

  error: string | null;
};

const FormAuth = (props: FormAuthProps): React.ReactNode => {
  const navigate = useNavigate();
  const form = useForm();

  const TITLE = props.isLogin ? "Login" : "Signup";
  const FOOTER_TEXT = props.isLogin
    ? "Don't have an account?"
    : "Already have an account?";
  const FOOTER_TEXT_LINK = props.isLogin ? "Create one" : "Login";

  return (
    <div className="flex justify-center items-center flex-col h-[calc(100svh)] dark:bg-black">
      <div className="flex flex-row w-full h-full">
        <div className="w-[50%] h-full xl:flex hidden flex-col items-center justify-center bg-accent/50">
          <div>
            <img src={MonsoonLogo} className="w-96 pointer-none mb-2" alt="" />
            <Label className="text-4xl">
              The safest way to communicate globally.
            </Label>
          </div>
        </div>
        <div className="xl:w-[50%] w-full h-full flex items-center justify-center">
          <div className="w-[50%] h-[50%]">
            <div className="flex justify-center items-center flex-col">
              <Card className="rounded-2xl p-5 border-none">
                <CardHeader>
                  <CardDescription className="w-full">
                    <Label htmlFor="login-title" className="text-2xl">
                      {TITLE}
                    </Label>
                  </CardDescription>
                </CardHeader>
                <CardContent className="w-96">
                  <Form {...form}>
                    <form onSubmit={props.onSubmit}>
                      <FormField
                        control={form.control}
                        name={TITLE.toLowerCase()}
                        render={() => (
                          <FormItem
                            className={`flex flex-col ${
                              props.isLogin ? "gap-3" : "gap-1"
                            }`}
                          >
                            {!props.isLogin && (
                              <>
                                <FormLabel className="flex flex-row items-center">
                                  <FiUser className="text-xl mx-2 text-primary-text" />
                                  Username
                                </FormLabel>
                                <FormControl>
                                  <Input
                                    className="rounded-xl"
                                    type="text"
                                    value={props.username}
                                    onChange={(e) =>
                                      props.setUsername!(e.target.value)
                                    }
                                  />
                                </FormControl>
                                <FormLabel className="flex flex-row items-center">
                                  <FaRegAddressCard className="text-xl mx-2 text-primary-text" />
                                  Display Name
                                </FormLabel>
                                <FormControl>
                                  <Input
                                    className="rounded-xl"
                                    type="text"
                                    value={props.displayName}
                                    onChange={(e) =>
                                      props.setDisplayName!(e.target.value)
                                    }
                                  />
                                </FormControl>
                              </>
                            )}
                            <FormLabel className="flex flex-row items-center">
                              <FiMail className="text-xl mx-2 text-primary-text" />{" "}
                              Email
                            </FormLabel>
                            <FormControl>
                              <Input
                                className="rounded-xl"
                                type="email"
                                value={props.email}
                                onChange={(e) => props.setEmail(e.target.value)}
                              />
                            </FormControl>
                            <FormLabel className="flex flex-row items-center">
                              <PiPasswordBold className="text-xl mx-2 text-primary-text" />{" "}
                              Password
                            </FormLabel>
                            <FormControl>
                              <Input
                                className="rounded-xl"
                                type="password"
                                value={props.password}
                                onChange={(e) =>
                                  props.setPassword(e.target.value)
                                }
                              />
                            </FormControl>
                            <div className="mt-10 w-full h-20 flex items-center">
                              <Button
                                type="submit"
                                className="w-full rounded-xl"
                              >
                                {TITLE}
                              </Button>
                            </div>
                          </FormItem>
                        )}
                      />
                    </form>
                  </Form>
                </CardContent>
                <CardFooter className="w-full flex justify-center">
                  <Label
                    htmlFor="signup"
                    className="text-white/30 text-center font-normal"
                  >
                    {FOOTER_TEXT}{" "}
                    <u>
                      <a
                        className="cursor-pointer hover:text-primary transition ease-linear"
                        onClick={() =>
                          navigate(props.isLogin ? "/register" : "/login")
                        }
                      >
                        {FOOTER_TEXT_LINK}
                      </a>
                    </u>
                  </Label>
                </CardFooter>
              </Card>
              <Label
                htmlFor="copyright"
                className="w-full text-center opacity-40 fixed bottom-[30px]"
              >
                © K.M Ahnaf Zamil, Monsoon Team 2025
              </Label>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const RegisterAuth = (): React.ReactNode => {
  const navigate = useNavigate();
  const currentUser = useCurrentUser();

  const [username, setUsername] = useState<string>("");
  const [displayName, setDisplayName] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (currentUser.isSuccess) {
      navigate("/");
    }
  }, [currentUser.isPending, navigate]);

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
          case 409:
            setError("User already exists");
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

  useEffect(() => {
    if (error) {
      toast({
        title: "Something went wrong",
        description: error,
        variant: "destructive",
      });
      setError(null);
    }
  }, [error]);

  if (currentUser.isPending) {
    return <LoadingPage />;
  } else if (currentUser.isSuccess) {
    return <Navigate to="/" />;
  } else {
    return (
      <FormAuth
        isLogin={false}
        onSubmit={handleSubmit}
        setEmail={setEmail}
        setPassword={setPassword}
        email={email}
        password={password}
        setDisplayName={setDisplayName}
        setUsername={setUsername}
        displayName={displayName}
        username={username}
        setError={setError}
        error={error}
      />
    );
  }
};

const LoginAuth = (): React.ReactNode => {
  const currentUser = useCurrentUser();

  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string | null>(null);

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
          decodeBase64(resp.data.enc_salt)
        );
        decryptedSeed = await CryptoHelper.AESGCMDecrypt(
          encKey,
          decodeBase64(resp.data.nonce),
          decodeBase64(resp.data.enc_seed)
        );
      } catch {
        return setError("Error while decrypting key seed");
      }

      await storeSeedSecure(decryptedSeed);

      window.location.href = "/";
    } catch (e) {
      console.error(e);
      setError("An error occured while logging you in");
    }
  };

  useEffect(() => {
    if (error) {
      toast({
        title: "Something went wrong",
        description: error,
        variant: "destructive",
      });
      setError(null);
    }
  }, [error]);

  if (currentUser.isPending) {
    return <LoadingPage />;
  } else if (currentUser.isSuccess) {
    return <Navigate to="/" />;
  } else {
    return (
      <FormAuth
        isLogin
        onSubmit={handleSubmit}
        setEmail={setEmail}
        setPassword={setPassword}
        email={email}
        password={password}
        setError={setError}
        error={error}
      />
    );
  }
};

export default function AuthForm({ isLogin }: { isLogin?: boolean }) {
  return <>{isLogin ? <LoginAuth /> : <RegisterAuth />}</>;
}
