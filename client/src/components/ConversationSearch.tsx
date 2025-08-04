import type React from "react";

export const ConversationSearch: React.FC = () => {
    return (
        <div className="flex absolute top-17 left-0 right-0 h-14 justify-center">
            <div className="flex items-center gap-5 w-full mx-5 relative">
                <div className="absolute h-4 pl-3 pointer-events-none text-neutral-600 dark:text-neutral-500"></div>
                <input
                    type="text"
                    className="pl-9 w-full rounded-md px-4 py-1 bg-neutral-200 placeholder:text-neutral-600 dark:bg-neutral-800 dark:placeholder:text-neutral-500 dark:text-white"
                    placeholder="Search"
                    value=""
                />
            </div>
        </div>
    );
};
