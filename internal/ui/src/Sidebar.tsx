import { Undo, Redo, RotateCw, Plus, Download } from "lucide-react";
import { pieces } from "./pieces";
import { useBoard } from "./BoardContext";
import { GetFEN, NewGame, PlayMove, UndoMove } from "../wailsjs/go/main/App";
import { useState } from "react";

function Sidebar() {
    const { state, setState, loadBoard } = useBoard();
    const [showExport, setShowExport] = useState(false);

    const toAlgebraic = (sq: number) => {
        const file = String.fromCharCode(97 + (sq % 8));
        const rank = Math.floor(sq / 8) + 1;
        return `${file}${rank}`;
    };

    async function handleUndo() {
        if (state.currentMoveIndex < 0) return;

        const entry = state.moveHistory[state.currentMoveIndex];
        await UndoMove(entry.move, entry.undo);

        setState((prev) => ({
            ...prev,
            currentMoveIndex: prev.currentMoveIndex - 1,
        }));
        await loadBoard();
    }

    async function handleRedo() {
        if (state.currentMoveIndex >= state.moveHistory.length - 1) return;

        const entry = state.moveHistory[state.currentMoveIndex + 1];
        const newUndo = await PlayMove(entry.move);

        setState((prev) => {
            const newHistory = [...prev.moveHistory];
            newHistory[prev.currentMoveIndex + 1] = { ...entry, undo: newUndo };
            return {
                ...prev,
                moveHistory: newHistory,
                currentMoveIndex: prev.currentMoveIndex + 1,
            };
        });
        await loadBoard();
    }

    async function handleJumpToMove(index: number) {
        const currentIndex = state.currentMoveIndex;

        if (index === currentIndex) return;

        if (index < currentIndex) {
            for (let i = currentIndex; i > index; i--) {
                const entry = state.moveHistory[i];
                await UndoMove(entry.move, entry.undo);
            }
        } else {
            for (let i = currentIndex + 1; i <= index; i++) {
                const entry = state.moveHistory[i];
                const newUndo = await PlayMove(entry.move);
                setState((prev) => {
                    const newHistory = [...prev.moveHistory];
                    newHistory[i] = { ...entry, undo: newUndo };
                    return { ...prev, moveHistory: newHistory };
                });
            }
        }

        setState((prev) => ({ ...prev, currentMoveIndex: index }));
        await loadBoard();
    }

    function handleFlipBoard() {
        setState((prev) => ({ ...prev, boardFlipped: !prev.boardFlipped }));
    }

    async function handleNewGame() {
        await NewGame();
        await loadBoard();
        setState((prev) => ({
            ...prev,
            moveHistory: [],
            currentMoveIndex: -1,
            marks: [],
            arrows: [],
            selectedSquare: null,
            legalMoves: [],
        }));
    }

    return (
        <div className="min-w-40 md:min-w-70 lg:min-w-80 text-white p-4 flex flex-col gap-4 h-screen">
            <div className="text-2xl font-bold text-center">Chess</div>
            <div className="flex-1 overflow-y-auto p-3">
                <div className="text-sm font-semibold mb-2">Move History</div>
                <div className="grid grid-cols-2 gap-1">
                    {state.moveHistory.map((move, idx) => {
                        const moveNumber = Math.floor(idx / 2) + 1;
                        const isWhite = idx % 2 === 0;
                        const isActive = idx === state.currentMoveIndex;

                        return (
                            <button
                                key={idx}
                                onClick={() => handleJumpToMove(idx)}
                                className={`px-2 py-1 rounded flex items-center gap-1 cursor-pointer ${
                                    isActive
                                        ? "bg-gray-900"
                                        : "hover:bg-gray-900"
                                } ${isWhite ? "col-start-1" : "col-start-2"}`}
                            >
                                <>
                                    <span className="text-gray-400 w-6">
                                        {moveNumber}.
                                    </span>
                                    <div className="w-5 h-5 flex items-center justify-center">
                                        {pieces[move.piece]?.()}
                                    </div>
                                </>
                                <span>
                                    {toAlgebraic(move.from)}-
                                    {toAlgebraic(move.to)}
                                </span>
                            </button>
                        );
                    })}
                </div>
            </div>
            <button
                onClick={() => setShowExport(true)}
                className="flex items-center justify-center gap-2 px-4 py-2 hover:bg-gray-800 rounded-lg transition-colors cursor-pointer"
            >
                <Download size={18} />
                <span>Export</span>
            </button>
            <div className="grid grid-cols-3 gap-2">
                <button
                    onClick={handleUndo}
                    disabled={state.currentMoveIndex < 0}
                    className="flex items-center justify-center p-3 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors cursor-pointer"
                >
                    <Undo size={20} />
                </button>
                <button
                    onClick={handleRedo}
                    disabled={
                        state.currentMoveIndex >= state.moveHistory.length - 1
                    }
                    className="flex items-center justify-center p-3 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors cursor-pointer"
                >
                    <Redo size={20} />
                </button>
                <button
                    onClick={handleFlipBoard}
                    className="flex items-center justify-center p-3 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors cursor-pointer"
                >
                    <RotateCw size={20} />
                </button>
            </div>
            <button
                onClick={handleNewGame}
                className="flex items-center justify-center gap-2 px-4 py-3 hover:bg-gray-800 rounded-lg transition-colors font-semibold cursor-pointer"
            >
                <Plus size={20} />
                <span>New Game</span>
            </button>
            {showExport && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-neutral-900 p-4 rounded-lg flex flex-col gap-2">
                        <button
                            onClick={async () => {
                                const fen = await GetFEN();
                                await navigator.clipboard.writeText(fen);
                                setShowExport(false);
                            }}
                            className="px-4 py-2 hover:bg-gray-800 rounded cursor-pointer"
                        >
                            Copy FEN
                        </button>
                        <button
                            onClick={() => setShowExport(false)}
                            className="px-4 py-2 hover:bg-gray-800 rounded cursor-pointer"
                        >
                            Cancel
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}

export default Sidebar;
