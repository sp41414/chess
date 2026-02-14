import { useBoard } from "./BoardContext";

function Arrows() {
    const { state } = useBoard();

    return (
        <svg
            viewBox="0 0 100 100"
            className="absolute top-0 left-0 pointer-events-none z-20"
            style={{ width: "100%", height: "100%" }}
            preserveAspectRatio="xMidYMid meet"
        >
            {state.arrows.map((arrow, i) => {
                const fromFile = arrow.from % 8;
                const fromRank = 7 - Math.floor(arrow.from / 8);
                const toFile = arrow.to % 8;
                const toRank = 7 - Math.floor(arrow.to / 8);

                const sqSize = 100 / 8;

                const fromX = (fromFile + 0.5) * sqSize;
                const fromY = (fromRank + 0.5) * sqSize;
                const toX = (toFile + 0.5) * sqSize;
                const toY = (toRank + 0.5) * sqSize;

                const dx = Math.abs(toFile - fromFile);
                const dy = Math.abs(toRank - fromRank);
                const isKnightMove =
                    (dx === 2 && dy === 1) || (dx === 1 && dy === 2);

                const shorten = sqSize * 0.4;

                let pathD: string;

                if (isKnightMove) {
                    const cornerX = dy > dx ? fromX : toX;
                    const cornerY = dy > dx ? toY : fromY;

                    const finalDx = toX - cornerX;
                    const finalDy = toY - cornerY;
                    const finalLen = Math.hypot(finalDx, finalDy);

                    const endX =
                        cornerX + (finalDx * (finalLen - shorten)) / finalLen;
                    const endY =
                        cornerY + (finalDy * (finalLen - shorten)) / finalLen;

                    pathD = `M${fromX},${fromY} L${cornerX},${cornerY} L${endX},${endY}`;
                } else {
                    const len = Math.hypot(toX - fromX, toY - fromY);
                    const endX =
                        fromX + ((toX - fromX) * (len - shorten)) / len;
                    const endY =
                        fromY + ((toY - fromY) * (len - shorten)) / len;

                    pathD = `M${fromX},${fromY} L${endX},${endY}`;
                }

                return (
                    <g key={`arrow-${i}`}>
                        <marker
                            id={`arrowhead-${i}`}
                            markerWidth="3"
                            markerHeight="3"
                            refX="1.5"
                            refY="1.5"
                            orient="auto"
                        >
                            <polygon points="0 0, 3 1.5, 0 3" fill="#15803d" />
                        </marker>
                        <path
                            d={pathD}
                            stroke="#15803d"
                            strokeWidth={sqSize / 6}
                            fill="none"
                            opacity="0.8"
                            markerEnd={`url(#arrowhead-${i})`}
                        />
                    </g>
                );
            })}
        </svg>
    );
}

export default Arrows;
