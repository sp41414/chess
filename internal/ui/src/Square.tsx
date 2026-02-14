import { useDraggable, useDroppable } from "@dnd-kit/core";
import { pieces } from "./pieces";

interface SquareProps {
    idx: number;
    isDark: boolean;
    piece: string | null;
    isSelected: boolean;
    isLegalMove: boolean;
    isKingInCheck: boolean;
    isCapture: boolean;
    isLastMove: boolean;
    isMarked: boolean;
    boardFlipped: boolean;
    canDrag: boolean;
    onSquareClick: (idx: number) => void;
    onRightMouseDown: (idx: number, e: React.MouseEvent) => void;
    onRightMouseUp: (idx: number, e: React.MouseEvent) => void;
}

function Square({
    idx,
    isDark,
    piece,
    isSelected,
    isLegalMove,
    isKingInCheck,
    isCapture,
    isLastMove,
    isMarked,
    boardFlipped,
    canDrag,
    onSquareClick,
    onRightMouseDown,
    onRightMouseUp,
}: SquareProps) {
    const { setNodeRef: setDropRef } = useDroppable({
        id: idx.toString(),
    });

    const {
        attributes,
        listeners,
        setNodeRef: setDragRef,
        transform,
        isDragging,
    } = useDraggable({
        id: idx.toString(),
        disabled: !canDrag || !piece,
    });

    const dragStyle = transform
        ? {
              transform: `translate3d(${transform.x}px, ${transform.y}px, 0) scale(1.1)`,
          }
        : undefined;

    const file = idx % 8;
    const rank = Math.floor(idx / 8);

    const isBottomRank = boardFlipped ? rank === 7 : rank === 0;
    const isLeftFile = boardFlipped ? file === 7 : file === 0;

    const fileLabel = String.fromCharCode(97 + file);
    const rankLabel = (rank + 1).toString();

    return (
        <div
            ref={setDropRef}
            onMouseDown={(e) => e.button === 2 && onRightMouseDown(idx, e)}
            onMouseUp={(e) => e.button === 2 && onRightMouseUp(idx, e)}
            onClick={() => onSquareClick(idx)}
            onContextMenu={(e) => e.preventDefault()}
            className={`relative flex items-center justify-center ${
                isLastMove
                    ? "bg-board-selected/80"
                    : isSelected
                      ? "bg-board-selected/90 ring-2 ring-inset border-none"
                      : isDark
                        ? "bg-board-dark"
                        : "bg-board-light"
            }`}
        >
            {isLeftFile && (
                <div
                    className={`absolute top-0.5 left-1 text-xs md:text-base lg:text-lg font-semibold ${
                        isDark ? "text-board-light" : "text-board-dark"
                    }`}
                >
                    {rankLabel}
                </div>
            )}
            {isBottomRank && (
                <div
                    className={`absolute bottom-0.5 right-1 text-xs md:text-base lg:text-lg font-semibold ${
                        isDark ? "text-board-light" : "text-board-dark"
                    }`}
                >
                    {fileLabel}
                </div>
            )}
            {isMarked && (
                <div className="absolute inset-0 bg-red-500/40 z-0 pointer-events-none" />
            )}
            {isKingInCheck && (
                <div className="absolute w-full h-full rounded-full bg-red-500/70 blur-sm z-0 will-change-transform" />
            )}
            {piece && (
                <div
                    ref={setDragRef}
                    style={dragStyle}
                    {...listeners}
                    {...attributes}
                    className={`${isDragging ? "z-50" : "z-10"} cursor-grab active:cursor-grabbing touch-none`}
                >
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
}

export default Square;
