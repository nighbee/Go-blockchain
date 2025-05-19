import React, { useState, useContext } from 'react';
import { WalletContext } from './WalletContext';
import { TransactionContext } from './TransactionContext';

const API_URL = 'http://localhost:5001';

function Mining() {
    const { wallets } = useContext(WalletContext);
    const { miningHistory, addMiningRecord } = useContext(TransactionContext);
    const [selectedWallet, setSelectedWallet] = useState('');
    const [miningStatus, setMiningStatus] = useState('');
    const [isMining, setIsMining] = useState(false);

    const handleMine = async () => {
        if (!selectedWallet) {
            setMiningStatus('Please select a wallet first');
            return;
        }

        setIsMining(true);
        setMiningStatus('Mining in progress...');

        try {
            const response = await fetch(`${API_URL}/mine`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    minerAddress: selectedWallet
                })
            });

            const data = await response.json();

            if (response.ok) {
                setMiningStatus('Mining successful!');
                addMiningRecord({
                    reward: data.reward,
                    status: 'success'
                });
            } else {
                setMiningStatus(`Mining failed: ${data.error}`);
                addMiningRecord({
                    error: data.error,
                    status: 'failed'
                });
            }
        } catch (error) {
            setMiningStatus(`Error: ${error.message}`);
            addMiningRecord({
                error: error.message,
                status: 'error'
            });
        } finally {
            setIsMining(false);
        }
    };

    return (
        <div style={{ margin: '20px' }}>
            <h2>Mining</h2>
            <div style={{ marginBottom: '20px' }}>
                <label style={{ display: 'block', marginBottom: '10px' }}>
                    Select Wallet:
                    <select
                        value={selectedWallet}
                        onChange={(e) => setSelectedWallet(e.target.value)}
                        style={{
                            width: '100%',
                            padding: '8px',
                            marginTop: '5px',
                            borderRadius: '4px'
                        }}
                    >
                        <option value="">Select a wallet</option>
                        {wallets.map((wallet, index) => (
                            <option key={index} value={wallet.address}>
                                Wallet {index + 1}: {wallet.address}
                            </option>
                        ))}
                    </select>
                </label>
                
                <button
                    onClick={handleMine}
                    disabled={isMining || !selectedWallet}
                    style={{
                        padding: '10px 20px',
                        backgroundColor: isMining ? '#ccc' : '#007bff',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: isMining ? 'not-allowed' : 'pointer',
                    }}
                >
                    {isMining ? 'Mining...' : 'Start Mining'}
                </button>
            </div>

            {miningStatus && (
                <div style={{
                    padding: '10px',
                    marginBottom: '20px',
                    backgroundColor: miningStatus.includes('successful') ? '#d4edda' : '#f8d7da',
                    borderRadius: '4px'
                }}>
                    {miningStatus}
                </div>
            )}

            <h3>Mining History</h3>
            {miningHistory.length === 0 ? (
                <p>No mining history yet.</p>
            ) : (
                <ul style={{ listStyle: 'none', padding: 0 }}>
                    {miningHistory.map((record, index) => (
                        <li
                            key={index}
                            style={{
                                border: '1px solid #ccc',
                                padding: '10px',
                                margin: '5px 0',
                                borderRadius: '5px',
                                backgroundColor: record.status === 'success' ? '#d4edda' : '#f8d7da'
                            }}
                        >
                            <p><strong>Time:</strong> {record.timestamp}</p>
                            {record.reward && <p><strong>Reward:</strong> {record.reward} coins</p>}
                            {record.error && <p><strong>Error:</strong> {record.error}</p>}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Mining;