// src/components/Wallet.js
import React, { useState, useContext } from 'react';
import { createWallet } from './api';
import { WalletContext } from './WalletContext';

function Wallet() {
    const { wallets, setWallets } = useContext(WalletContext);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleCreateWallet = async () => {
        setLoading(true);
        setError(null);
        try {
            const wallet = await createWallet();
            setWallets([...wallets, wallet]);
        } catch (err) {
            setError('Failed to create wallet');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ margin: '20px' }}>
            <h2>Wallets</h2>
            <button
                onClick={handleCreateWallet}
                disabled={loading}
                style={{
                    padding: '10px 20px',
                    backgroundColor: '#007bff',
                    color: 'white',
                    border: 'none',
                    borderRadius: '5px',
                    cursor: loading ? 'not-allowed' : 'pointer',
                }}
            >
                {loading ? 'Creating...' : 'Create New Wallet'}
            </button>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            <h3>Registered Wallets</h3>
            {wallets.length === 0 ? (
                <p>No wallets created yet.</p>
            ) : (
                <ul style={{ listStyle: 'none', padding: 0 }}>
                    {wallets.map((wallet, index) => (
                        <li
                            key={index}
                            style={{
                                border: '1px solid #ccc',
                                padding: '10px',
                                margin: '5px 0',
                                borderRadius: '5px',
                            }}
                        >
                            <p><strong>Wallet {index + 1} Address:</strong> {wallet.address}</p>
                            <p><strong>Status:</strong> {wallet.message}</p>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Wallet;