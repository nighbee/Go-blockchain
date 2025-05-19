// src/App.js
import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { WalletProvider } from './components/WalletContext';
import { TransactionProvider } from './components/TransactionContext';
import Wallet from './components/Wallet';
import Transactions from './components/Transactions';
import Blockchain from './components/Blockchain';
import Mining from './components/Mining';
import Nodes from './components/Nodes';
import './App.css';

function App() {
    return (
        <WalletProvider>
            <TransactionProvider>
                <Router>
                    <div>
                        <h1>Blockchain Explorer</h1>
                        <nav style={{ padding: '20px', backgroundColor: '#f8f9fa' }}>
                            <ul style={{ 
                                listStyle: 'none', 
                                display: 'flex', 
                                gap: '20px',
                                margin: 0,
                                padding: 0
                            }}>
                                <li><Link to="/wallet">Wallet</Link></li>
                                <li><Link to="/transactions">Transactions</Link></li>
                                <li><Link to="/blockchain">Blockchain</Link></li>
                                <li><Link to="/mining">Mining</Link></li>
                                <li><Link to="/nodes">Nodes</Link></li>
                            </ul>
                        </nav>

                        <Routes>
                            <Route path="/wallet" element={<Wallet />} />
                            <Route path="/transactions" element={<Transactions />} />
                            <Route path="/blockchain" element={<Blockchain />} />
                            <Route path="/mining" element={<Mining />} />
                            <Route path="/nodes" element={<Nodes />} />
                            <Route path="/" element={<Wallet />} />
                        </Routes>
                    </div>
                </Router>
            </TransactionProvider>
        </WalletProvider>
    );
}

export default App;