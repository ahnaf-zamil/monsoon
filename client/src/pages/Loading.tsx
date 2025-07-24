import type React from "react";
import MonsoonLogo from "../static/img/monsoon_logo.png";
import { appVersion } from "../utils";

export const LoadingPage: React.FC = () => {
    return (
        <>
            <div className="h-screen w-full flex items-center justify-center  gap-10 flex-col">
                <img
                    src={MonsoonLogo}
                    className="w-52 md:w-96 pointer-none"
                    alt="Monsoon"
                    draggable="false"
                />
                <div className="animate-spin md:h-15 md:w-15 h-5 w-5 md:border-l-5 md:border-r-5 border-l-2 border-r-2 border-secondary-text rounded-full"></div>
            </div>
            <p className="fixed font-light bottom-5 text-sm text-center w-full  text-slate-500 pt-4">
                Version {appVersion}
            </p>
        </>
    );
};
