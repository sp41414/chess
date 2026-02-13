import { pieces } from "./pieces";

interface PromotionProps {
    isWhite: boolean;
    onSelect: (promotionFlag: number) => void;
}

function Promotion({ isWhite, onSelect }: PromotionProps) {
    const promotionPieces = [
        { piece: "Q", flag: 11 },
        { piece: "R", flag: 10 },
        { piece: "B", flag: 9 },
        { piece: "N", flag: 8 },
    ];

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
            <div className="bg-neutral-900 rounded-lg p-6 flex gap-4 shadow-2xl border border-gray-700">
                {promotionPieces.map(({ piece, flag }) => {
                    const pieceKey = `${isWhite ? "w" : "b"}${piece}`;

                    return (
                        <button
                            key={piece}
                            onClick={() => onSelect(flag)}
                            className="w-24 h-24 hover:bg-gray-800 rounded-lg flex items-center justify-center transition-colors"
                        >
                            {pieces[pieceKey]?.()}
                        </button>
                    );
                })}
            </div>
        </div>
    );
}

export default Promotion;
