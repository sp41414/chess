import { createContext } from "react";
import type { BoardContextType } from "./types";

export const BoardContext = createContext<BoardContextType | null>(null);
