import React from "react";
import { Home } from "./pages/Home";
import { Route, Routes } from "react-router-dom";
import { Login } from "./pages/Login";
import { AuthRequired } from "./wrappers/AuthRequired";
import { Register } from "./pages/Register";

const App: React.FC = () => (
    <Routes>
        <Route element={<AuthRequired />}>
            <Route path="/" element={<Home />} />
            <Route path="/conversations/:conversationID" element={<Home/>} />
        </Route>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
    </Routes>
);

export default App;
