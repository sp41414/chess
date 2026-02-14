import Board from "./Board";
import Sidebar from "./Sidebar";

function App() {
    return (
        <div className="bg-gray-950 min-h-screen flex justify-around">
            <div className="flex-1 flex items-center justify-center">
                <Board />
            </div>
            <Sidebar />
        </div>
    );
}

export default App;
