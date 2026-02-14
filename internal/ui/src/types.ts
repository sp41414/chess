import { engine } from "../wailsjs/go/models";

export type Arrow = {
    from: number;
    to: number;
};

export type MoveHistoryEntry = {
    from: number;
    to: number;
    piece: string;
    move: number;
    undo: engine.Undo;
};

export type BoardState = {
    pieces: Record<SquareIndex, string | null>;
    selectedSquare: SquareIndex | null;
    legalMoves: SquareIndex[];
    arrows: Arrow[];
    marks: SquareIndex[];
    moveHistory: MoveHistoryEntry[];
    currentMoveIndex: number;
    boardFlipped: boolean;
};

export type BoardContextType = {
    state: BoardState;
    setState: React.Dispatch<React.SetStateAction<BoardState>>;
    loadBoard: () => Promise<void>;
};

export type PieceRender = Record<
    string,
    (props?: {
        fill?: string;
        square?: string;
        svgStyle?: React.CSSProperties;
    }) => React.JSX.Element
>;

// 0-63, a1 is 0, h8 is 63
export type SquareIndex = number;

export type Move = {
    from: number;
    to: number;
    flags?: number;
};

export type FenString =
    | "p"
    | "r"
    | "n"
    | "b"
    | "q"
    | "k"
    | "P"
    | "R"
    | "N"
    | "B"
    | "Q"
    | "K";
