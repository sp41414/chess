import { useEffect, useState } from "react";
import type { BoardState } from "./types";
import { BoardContext } from "./BoardContext";
import { GetPieces } from "../wailsjs/go/main/App";

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

    async function loadBoard() {
        const pieces = await GetPieces();
        setState((prev) => ({ ...prev, pieces }));
    }

    useEffect(() => {
        GetPieces().then((pieces) => {
            setState((prev) => ({ ...prev, pieces }));
        });
    }, []);

    return (
        <BoardContext.Provider value={{ state, setState, loadBoard }}>
            {children}
        </BoardContext.Provider>
    );
}

export default BoardProvider;
