import { useState } from "react";
import type { BoardState } from "./types";
import { BoardContext } from "./BoardContext";

function BoardProvider({ children }: { children: React.ReactNode }) {
    const [state, setState] = useState<BoardState>({
        // Square indexes 0-63
        pieces: {},
        selectedSquare: null,
        legalMoves: [],
        // Right click and hold arrows
        arrows: [],
        // Right click marks
        marks: [],
    });

    return (
        <BoardContext.Provider value={{ state, setState }}>
            {children}
        </BoardContext.Provider>
    );
}

export default BoardProvider;
