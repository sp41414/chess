import { createContext, useContext } from "react";
import type { BoardContextType } from "./types";

export const BoardContext = createContext<BoardContextType | null>(null);

export function useBoard() {
    const context = useContext(BoardContext);
    if (!context) {
        throw new Error("useBoard must be used within a BoardProvider");
    }
    return context;
}
