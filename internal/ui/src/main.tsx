import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import BoardProvider from "./BoardProvider";
import App from "./App";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <BoardProvider>
            <App />
        </BoardProvider>
    </StrictMode>,
);
