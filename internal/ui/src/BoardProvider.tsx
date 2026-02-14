import { useEffect, useState } from "react";
import type { BoardState } from "./types";
import { BoardContext } from "./BoardContext";
import { GetPieces } from "../wailsjs/go/main/App";

function BoardProvider({ children }: { children: React.ReactNode }) {
    const [state, setState] = useState<BoardState>({
        pieces: {},
        selectedSquare: null,
        legalMoves: [],
        arrows: [],
        marks: [],
        moveHistory: [],
        currentMoveIndex: -1,
        boardFlipped: false,
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
