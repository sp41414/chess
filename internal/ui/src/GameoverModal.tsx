type GameoverModalProps = {
    winner: string | null;
    gameoverType: string;
    onNewGame: () => void;
    onClose: () => void;
};

function GameoverModal({
    winner,
    gameoverType,
    onNewGame,
    onClose,
}: GameoverModalProps) {
    return (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
            <div className="bg-neutral-800 rounded-xl p-8 max-w-md w-full mx-4 shadow-2xl border border-neutral-700">
                <div className="text-center space-y-4">
                    <div className="text-3xl font-bold text-white">
                        {winner
                            ? `${winner === "w" ? "White" : "Black"} wins`
                            : "Draw"}
                    </div>
                    <div className="text-lg text-gray-400 capitalize">
                        by {gameoverType}
                    </div>
                    <div className="flex gap-3">
                        <button
                            onClick={onClose}
                            className="flex-1 px-6 py-3 bg-neutral-700 hover:bg-neutral-600 text-white font-semibold rounded-lg transition-colors cursor-pointer"
                        >
                            Review
                        </button>
                        <button
                            onClick={onNewGame}
                            className="flex-1 px-6 py-3 bg-green-600 hover:bg-green-700 text-white font-semibold rounded-lg transition-colors cursor-pointer"
                        >
                            New Game
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default GameoverModal;
