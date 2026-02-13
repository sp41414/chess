export type Arrow = {
    startSquare: string;
    endSquare: string;
    color: string;
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

export type BoardState = {
    pieces: Record<SquareIndex, string | null>;
    selectedSquare: SquareIndex | null;
    legalMoves: SquareIndex[];

    arrows: Arrow[];
    marks: SquareIndex[];
};

export type BoardContextType = {
    state: BoardState;
    setState: React.Dispatch<React.SetStateAction<BoardState>>;
    loadBoard: () => Promise<void>;
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
