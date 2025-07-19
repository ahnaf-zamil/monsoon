import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import * as path from "path";

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            "argon2-browser": path.resolve(
                __dirname,
                "node_modules/argon2-browser/dist/argon2-bundled.min.js"
            ),
        },
    },
});
