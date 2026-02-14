import { useEffect, useState } from "react";
import type { BoardState } from "./types";
import { BoardContext } from "./BoardContext";
import { GetFEN, GetPieces } from "../wailsjs/go/main/App";

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
        sideToMove: "w",
    });

    async function loadBoard() {
        const pieces = await GetPieces();
        const fen = await GetFEN();
        const sideToMove = fen.split(" ")[1] as "w" | "b";
        setState((prev) => ({ ...prev, pieces, sideToMove }));
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
