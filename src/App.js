// src/App.js
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Blockchain from './components/Blockchain';
import Wallet from './components/Wallet';
import Transactions from './components/Transactions';
import { WalletProvider } from './components/WalletContext'; // Correct path
import './App.css';

function App() {
    return (
        <WalletProvider>
            <Router>
                <div>
                    <h1>Blockchain Explorer</h1>
                    <nav style={{ marginBottom: '20px' }}>
                        <ul style={{ listStyle: 'none', padding: 0, display: 'flex', gap: '20px' }}>
                            <li>
                                <Link to="/" style={{ textDecoration: 'none', color: '#007bff' }}>
                                    Blockchain
                                </Link>
                            </li>
                            <li>
                                <Link to="/wallet" style={{ textDecoration: 'none', color: '#007bff' }}>
                                    Wallets
                                </Link>
                            </li>
                            <li>
                                <Link to="/transactions" style={{ textDecoration: 'none', color: '#007bff' }}>
                                    Transactions
                                </Link>
                            </li>
                        </ul>
                    </nav>
                    <Routes>
                        <Route path="/" element={<Blockchain />} />
                        <Route path="/wallet" element={<Wallet />} />
                        <Route path="/transactions" element={<Transactions />} />
                    </Routes>
                </div>
            </Router>
        </WalletProvider>
    );
}

export default App;