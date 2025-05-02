import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Blockchain from './components/Blockchain';
import Nodes from './components/Nodes';
import Wallet from './components/Wallet';
import './App.css';

function App() {
    return (
        <Router>
            <div>
                <h1>Blockchain Explorer</h1>
                <Routes>
                    <Route path="/" element={<Blockchain />} />
                    <Route path="/nodes" element={<Nodes />} />
                    <Route path="/wallet" element={<Wallet />} />
                </Routes>
            </div>
        </Router>
    );
}

export default App;