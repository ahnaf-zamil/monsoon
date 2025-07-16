import type React from "react";
import { useState } from "react";
import { BiCheck, BiImageAlt } from "react-icons/bi";
import { FiEdit2 } from "react-icons/fi";
import { isEmptyString } from "../util";

interface MessageBoxProps {
  submitHandler: (content: string) => void;
}

export const MessageBox: React.FC<MessageBoxProps> = ({ submitHandler }) => {
  const isEditing = false;
  const [content, setContent] = useState<string>("");

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!isEmptyString(content)) {
      // Non-empty and
      submitHandler(content);
      setContent("");
    }
  };

  return (
    <>
      <form
        onSubmit={(e) => {
          onSubmit(e);
        }}
        className="absolute bottom-0 flex-grow h-20 px-5 flex items-center gap-2 w-full"
      >
        {/* {imgBase64 && (
          <div className="absolute top-0 -translate-y-full">
            <div className="relative w-24">
              <button
                type="button"
                className="bg-neutral-400 absolute flex items-center rounded-full aspect-square p-2 top-0 right-0 translate-x-1/2 -translate-y-1/2 cursor-pointer"
                onClick={() => {
                  setImgBase64("");
                }}
              >
                <IoClose size="1rem" />
              </button>
              <img
                src={imgBase64}
                alt="message image"
                className="aspect-square object-cover rounded-xl"
              />
            </div>
          </div>
        )} */}
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
              className={`outline outline-2 outline-primary w-full rounded-full px-6 py-3 bg-neutral-200 placeholder:text-neutral-600 dark:bg-neutral-800 dark:placeholder:text-neutral-500 dark:text-white`}
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
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
            type="submit"
            className={`bg-neutral-200 rounded-full h-12 aspect-square flex items-center justify-center p-2.5 text-primary dark:bg-neutral-800`}
          >
            <BiCheck size={"100%"} />
          </button>
        </div>
      </form>
    </>
  );
};
