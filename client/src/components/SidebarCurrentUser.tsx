import type React from "react";
import type { IUser } from "../types";
import { IoMdLogOut } from "react-icons/io";

interface Props {
    currentUser?: IUser;
}

export const SidebarCurrentUser: React.FC<Props> = ({ currentUser }) => {
    return (
        <div className="flex items-center justify-between px-5 py-3 gap-5 rounded-lg bg-zinc-900">
            <div className="flex items-center gap-2 hover:bg-zinc-800 hover:cursor-pointer rounded-l-full w-full rounded-r-full">
                <div className="w-12 aspect-square rounded-full overflow-hidden">
                    <img
                        src={undefined}
                        alt=""
                        className="object-cover w-full h-full bg-primary"
                    />
                </div>
                <div className="grid items-center">
                    <>
                        <p className="text-base leading-4 dark:text-white">
                            {currentUser?.display_name}
                        </p>
                        <p className="text-neutral-600 text-sm leading-4 dark:text-neutral-500">
                            @{currentUser?.username}
                        </p>
                    </>
                </div>
            </div>
            <button
                name="logout"
                className="h-fit p-3 rounded-full dark:hover:bg-neutral-800 dark:bg-neutral-900 dark:text-white hover:cursor-pointer"
            >
                <IoMdLogOut />
            </button>
        </div>
    );
};
