import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { BrowserRouter } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext.tsx";
import { SocketProvider } from "./context/SocketContext.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <BrowserRouter>
            <QueryClientProvider client={queryClient}>
                <ReactQueryDevtools
                    initialIsOpen={false}
                    buttonPosition="top-right"
                />
                <AuthProvider>
                    <SocketProvider>
                        <App />
                    </SocketProvider>
                </AuthProvider>
            </QueryClientProvider>
        </BrowserRouter>
    </StrictMode>,
);
