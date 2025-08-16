import type React from "react";
import type { IUser } from "@/types";
import { IoMdLogOut } from "react-icons/io";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";

export const SidebarCurrentUser = ({
  currentUser,
}: {
  currentUser?: IUser;
}): React.ReactNode => {
  return (
    <div className="flex items-center justify-between p-3 gap-5 rounded-lg bg-foreground/5 border-t w-full">
      <div className="flex items-center gap-2 hover:bg-foreground/10 hover:cursor-pointer transition ease-in w-full rounded-xl p-3">
        <div className="flex flex-row gap-2 items-center justify-center">
          <div className="w-12 aspect-square rounded-full overflow-hidden">
            <Avatar className="w-14 h-14">
              <AvatarImage src="https://github.com/ahnaf-zamil.png" alt="" />
              <AvatarFallback>AZ</AvatarFallback>
            </Avatar>
          </div>
          <div className="flex flex-col gap-1 items-start justify-center">
            <>
              <Label className="hover:cursor-pointer">
                {currentUser?.display_name}
              </Label>
              <Label className="text-[0.8rem] text-foreground/30 font-normal hover:cursor-pointer">
                @{currentUser?.username}
              </Label>
            </>
          </div>
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
