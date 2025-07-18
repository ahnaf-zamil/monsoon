import React from "react";
import { Home } from "./pages/Home";
import { Route, Routes } from "react-router-dom";
import { Login } from "./pages/Login";
import { AuthRequired } from "./pages/AuthRequired";

const App: React.FC = () => (
    <Routes>
        <Route element={<AuthRequired />}>
            <Route path="/" element={<Home />} />
        </Route>
        <Route path="/login" element={<Login />} />
    </Routes>
);

export default App;
