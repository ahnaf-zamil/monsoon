import type React from "react";
import MonsoonLogo from "../static/img/monsoon_logo.png";
import { useContext, useEffect, useState } from "react";
import { isEmptyString } from "../util";
import { loginUser } from "../api/auth";
import { isAxiosError } from "axios";
import { BiUser } from "react-icons/bi";
import { PiPasswordBold } from "react-icons/pi";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";

export const Login: React.FC = () => {
    const navigate = useNavigate();
    const currentUser = useContext(AuthContext);

  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (currentUser) {
            navigate("/");
        }
    }, [currentUser])

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!isEmptyString(username) && !isEmptyString(password)) {
      // Login user here
      const resp = await loginUser(username, password);
      if (!isAxiosError(resp)) {
        navigate("/");
      } else {
        switch (resp.response?.status) {
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
              <label htmlFor="username" className="sr-only">
                Username
              </label>
              <div className="flex w-full rounded-md dark:bg-neutral-800 items-center justify-center">
                <BiUser className="text-xl mx-2 text-primary-text" />
                <input
                  type="text"
                  className="outline-none rounded-md flex-grow pr-4 py-2 placeholder:text-neutral-600 dark:placeholder:text-neutral-500 dark:text-white"
                  placeholder="Username"
                  id="username"
                  required={true}
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
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
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>
            </div>
            <button className="hover:cursor-pointer hover:bg-primary-darker bg-primary p-2 rounded-md  flex justify-center text-black">
              Log in
            </button>
          </div>
        </form>
      </div>
      <p className="text-center text-gray-400 fixed bottom-[30px] ">
        © K.M Ahnaf Zamil 2025
      </p>
    </div>
  );
};
