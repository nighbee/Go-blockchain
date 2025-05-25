import React, { useState, useContext, useEffect } from 'react';
import { WalletContext } from './WalletContext';
import { TransactionContext } from './TransactionContext';
import { getBlockchain } from './api';

const API_URL = 'http://localhost:5001';

function Mining() {
    const { wallets } = useContext(WalletContext);
    const { addMiningRecord } = useContext(TransactionContext);
    const [selectedWallet, setSelectedWallet] = useState('');
    const [miningStatus, setMiningStatus] = useState('');
    const [isMining, setIsMining] = useState(false);
    const [miningRewards, setMiningRewards] = useState([]);

    // Fetch mining rewards from blockchain
    useEffect(() => {
        const fetchMiningRewards = async () => {
            try {
                const data = await getBlockchain();
                if (data.chain && data.chain.length > 0) {
                    // Collect all mining reward transactions
                    const rewards = data.chain.flatMap(block => 
                        (block.transactions || [])
                            .filter(tx => tx.message === "MINING REWARD")
                            .map(tx => ({
                                ...tx,
                                timestamp: new Date(block.timestamp * 1000).toLocaleString(),
                                blockIndex: block.index,
                                blockHash: block.hash
                            }))
                    );
                    setMiningRewards(rewards);
                }
            } catch (err) {
                console.error('Error fetching mining rewards:', err);
            }
        };
        fetchMiningRewards();
    }, []);

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
                // Refresh mining rewards after successful mining
                const blockchainData = await getBlockchain();
                if (blockchainData.chain && blockchainData.chain.length > 0) {
                    const latestBlock = blockchainData.chain[blockchainData.chain.length - 1];
                    const newReward = latestBlock.transactions.find(tx => tx.message === "MINING REWARD");
                    if (newReward) {
                        setMiningRewards(prev => [{
                            ...newReward,
                            timestamp: new Date(latestBlock.timestamp * 1000).toLocaleString(),
                            blockIndex: latestBlock.index,
                            blockHash: latestBlock.hash
                        }, ...prev]);
                    }
                }
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
                    Select Wallet to Mine:
                    <select
                        value={selectedWallet}
                        onChange={(e) => setSelectedWallet(e.target.value)}
                        style={{ width: '100%', padding: '8px', marginTop: '5px' }}
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
                        backgroundColor: '#007bff',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: isMining || !selectedWallet ? 'not-allowed' : 'pointer',
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
            {miningRewards.length === 0 ? (
                <p>No mining rewards yet.</p>
            ) : (
                <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
                    {miningRewards.map((reward, index) => (
                        <div
                            key={index}
                            style={{
                                border: '1px solid #ddd',
                                padding: '15px',
                                marginBottom: '10px',
                                borderRadius: '5px',
                                backgroundColor: '#f9f9f9'
                            }}
                        >
                            <p><strong>Time:</strong> {reward.timestamp}</p>
                            <p><strong>Block:</strong> #{reward.blockIndex}</p>
                            <p><strong>Miner:</strong> {reward.recipientBlockchainAddress}</p>
                            <p><strong>Reward:</strong> {reward.value} coins</p>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

export default Mining;