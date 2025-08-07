import type React from "react";
import { useState } from "react";
import { BiCheck, BiImageAlt } from "react-icons/bi";
import { FiEdit2 } from "react-icons/fi";
import { isEmptyString } from "../utils";
import { useInboxStore } from "../store/inbox";

interface MessageBoxProps {
    submitHandler: (content: string) => void;
}

export const MessageBox: React.FC<MessageBoxProps> = ({ submitHandler }) => {
    const inboxStore = useInboxStore();
    const selectedConversation = inboxStore.getSelectedConversation();

    const isEditing = false;
    const [content, setContent] = useState<string>("");

    const onSubmit = (e: React.SyntheticEvent) => {
        e.preventDefault();
        if (!isEmptyString(content)) {
            submitHandler(content);
            setContent("");
        }
    };

    return (
        <>
            <div className="bg-chatbox absolute bottom-0 flex-grow h-20 px-5 flex items-center gap-2 w-full">
                <div className="flex items-center gap-2 w-full my-3">
                    <div className="w-full relative flex flex-col">
                        {isEditing && (
                            <div className="absolute top-1/2 -translate-y-1/2 pl-3 text-blue-600 flex items-center gap-2 border-r-[1px] pr-2 border-blue-600">
                                <FiEdit2 />
                                <p className="text-sm">Editing</p>
                            </div>
                        )}
                        <input
                            type="text"
                            className={`border-1 border-zinc-700 focus:border-zinc-600 outline-none rounded-md w-full px-6 py-3 bg-neutral-200 placeholder:text-neutral-600 dark:bg-neutral-800 dark:placeholder:text-neutral-500 dark:text-white`}
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            onKeyDown={(e) => {
                                if (e.key == "Enter") onSubmit(e);
                            }}
                            placeholder={
                                "Message " +
                                (selectedConversation
                                    ? selectedConversation.name
                                    : "")
                            }
                        />
                    </div>

                    {!isEditing && (
                        <>
                            <input
                                type="file"
                                accept="image/jpeg, image/png"
                                id="file"
                                className="hidden"
                            />
                            <button
                                onClick={(e) => {
                                    e.preventDefault();
                                    console.log("CLICK");
                                    document.getElementById("file")?.click();
                                }}
                                type="button"
                                className="bg-neutral-200 rounded-full h-12 aspect-square flex items-center justify-center p-2.5 text-primary dark:bg-neutral-800"
                                aria-label="Upload Image"
                            >
                                <BiImageAlt size={"100%"} />
                            </button>
                        </>
                    )}

                    <button
                        onClick={(e) => onSubmit(e)}
                        className={`bg-neutral-200 rounded-full h-12 aspect-square flex items-center justify-center p-2.5 text-primary dark:bg-neutral-800`}
                    >
                        <BiCheck size={"100%"} />
                    </button>
                </div>
            </div>
        </>
    );
};
