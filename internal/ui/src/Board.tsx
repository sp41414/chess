import { useState } from "react";
import {
    GetFEN,
    GetMoves,
    GetPieces,
    IsInCheck,
    PlayMove,
} from "../wailsjs/go/main/App";
import { useBoard } from "./BoardContext";
import { pieces } from "./pieces";
import type { Move } from "./types";
import Promotion from "./Promotion";

function Board() {
    const RANKS = [7, 6, 5, 4, 3, 2, 1, 0];
    const FILES = [0, 1, 2, 3, 4, 5, 6, 7];

    const { state, setState, loadBoard } = useBoard();
    const [availableMoves, setAvailableMoves] = useState<number[]>([]);
    const [kingInCheckSq, setKingInCheckSq] = useState<number | null>(null);
    const [lastMove, setLastMove] = useState<Move | null>(null);
    const [promotion, setPromotion] = useState<Move | null>(null);

    const getMove = (move: number) => ({
        from: move & 0x3f,
        to: (move >> 6) & 0x3f,
        flags: (move >> 12) & 0xf,
    });

    async function findKingSq() {
        const check = await IsInCheck();
        if (!check) {
            setKingInCheckSq(null);
            return;
        }

        const pieces = await GetPieces();
        const fen = await GetFEN();
        const king = fen.split(" ")[1] === "w" ? "wK" : "bK";

        for (let sq = 0; sq < 64; sq++) {
            if (pieces[sq] === king) {
                console.log("king in check", sq);
                setKingInCheckSq(sq);
                return;
            }
        }
    }

    async function selectPiece(sq: number) {
        if (!state.pieces[sq]) return;
        const moves = await GetMoves();

        const movesFrom = moves.filter((m: number) => getMove(m).from === sq);

        setAvailableMoves(movesFrom);
        setState((prev) => ({
            ...prev,
            selectedSquare: sq,
            legalMoves: movesFrom.map((m: number) => getMove(m).to),
        }));
    }

    function clearSelection() {
        setAvailableMoves([]);
        setState((prev) => ({
            ...prev,
            selectedSquare: null,
            legalMoves: [],
        }));
    }

    async function handleMove(to: number) {
        const from = state.selectedSquare!;
        const move = availableMoves.find((m: number) => {
            const { from: fromSq, to: toSq } = getMove(m);
            return from === fromSq && to === toSq;
        });

        if (move) {
            const { flags } = getMove(move);
            // Check if the flag is a promotion
            if (flags >= 8) {
                setPromotion({ from, to });
                clearSelection();
                return;
            }

            setLastMove({ from, to });
            clearSelection();

            await PlayMove(move);
            await loadBoard();
            await findKingSq();
        } else if (state.pieces[to]) {
            clearSelection();
            await selectPiece(to);
        } else {
            clearSelection();
        }
    }

    async function handleSquareClick(sq: number) {
        if (state.selectedSquare === null) {
            await selectPiece(sq);
        } else {
            await handleMove(sq);
        }
    }

    async function handlePromotion(promotionFlag: number) {
        if (!promotion) return;

        const { from, to } = promotion;
        const isCapture = state.pieces[to] !== undefined;
        const finalFlag = isCapture ? promotionFlag + 4 : promotionFlag;

        const promotionMove =
            (from & 0x3f) | ((to & 0x3f) << 6) | ((finalFlag & 0xf) << 12);

        setLastMove({ from, to });
        setPromotion(null);

        await PlayMove(promotionMove);
        await loadBoard();
        await findKingSq();
    }

    return (
        <>
            <div className="grid grid-cols-8 grid-rows-8 aspect-square select-none max-w-[min(90vw,90vh)] max-h-[min(90vw,90vh)]">
                {RANKS.map((rank) =>
                    FILES.map((file) => {
                        const idx = rank * 8 + file;
                        const isDark = (rank + file) % 2 === 0;
                        const piece = state.pieces[idx];
                        const isSelected = state.selectedSquare === idx;
                        const isLegalMove = state.legalMoves.includes(idx);
                        const isKingInCheck = kingInCheckSq === idx;
                        const isCapture = isLegalMove && piece;
                        const isLastMove =
                            lastMove?.from === idx || lastMove?.to === idx;

                        return (
                            <div
                                key={idx}
                                onClick={() => handleSquareClick(idx)}
                                className={`relative flex items-center justify-center ${
                                    isLastMove
                                        ? "bg-board-selected/80"
                                        : isSelected
                                          ? "bg-board-selected/90 ring-2 ring-inset"
                                          : isDark
                                            ? "bg-board-dark"
                                            : "bg-board-light"
                                }`}
                            >
                                {isKingInCheck && (
                                    <div className="absolute w-full h-full rounded-full bg-red-500/70 blur-sm z-0 will-change-transform" />
                                )}
                                {piece && (
                                    <div className="z-10">
                                        {pieces[piece]?.()}
                                    </div>
                                )}
                                {isLegalMove && !isCapture && (
                                    <div className="absolute w-1/3 h-1/3 bg-black/30 rounded-full pointer-events-none" />
                                )}
                                {isCapture && (
                                    <div className="absolute w-4/5 h-4/5 ring-8 ring-red-400 rounded-full pointer-events-none z-0" />
                                )}
                            </div>
                        );
                    }),
                )}
            </div>
            {promotion && (
                <Promotion
                    isWhite={state.pieces[promotion.from]?.[0] === "w"}
                    onSelect={handlePromotion}
                />
            )}
        </>
    );
}

export default Board;
