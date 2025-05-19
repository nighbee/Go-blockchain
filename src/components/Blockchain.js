import React, { useEffect, useState } from 'react';
import { getBlockchain } from './api';

function Blockchain() {
    const [blocks, setBlocks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchBlockchain = async () => {
            try {
                const data = await getBlockchain();
                console.log('Fetched blockchain data:', data);
                if (data.chain && data.chain.length > 0) {
                    // Map backend fields to frontend expectations
                    const mappedBlocks = data.chain.map((block, index) => ({
                        index: index,
                        timestamp: block.timestamp * 1000, // Convert seconds to milliseconds
                        hash: block.hash || block.Hash,
                        previousHash: block.prevHash,
                        nonce: block.nonce,
                        transactions: (block.transactions || []).map(tx => ({
                            from: tx.senderBlockchainAddress,
                            to: tx.recipientBlockchainAddress,
                            amount: tx.value,
                            message: tx.message,
                        })),
                    }));
                    setBlocks(mappedBlocks);
                } else {
                    setError('Blockchain data is empty');
                }
            } catch (err) {
                console.error('Error fetching blockchain:', err);
                setError('Failed to load blockchain data');
            } finally {
                setLoading(false);
            }
        };
        fetchBlockchain();
    }, []);

    if (loading) {
        return <p>Loading blockchain data...</p>;
    }

    if (error) {
        return <p>{error}</p>;
    }

    return (
        <div style={{ margin: '20px' }}>
            <h2>Blockchain Activity</h2>
            {blocks.map((block) => (
                <div
                    key={block.index}
                    style={{
                        border: '1px solid #ddd',
                        margin: '10px 0',
                        padding: '15px',
                        borderRadius: '8px',
                        backgroundColor: '#f8f9fa',
                        boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
                    }}
                >
                    <div style={{ marginBottom: '10px' }}>
                        <h3 style={{ margin: '0 0 10px 0', color: '#2c3e50' }}>
                            Block #{block.index}
                        </h3>
                        <div style={{
                            display: 'grid',
                            gridTemplateColumns: '1fr 1fr',
                            gap: '10px',
                            fontSize: '0.9em'
                        }}>
                            <div>
                                <p><strong>Timestamp:</strong> {new Date(block.timestamp).toLocaleString()}</p>
                                <p><strong>Nonce:</strong> {block.nonce}</p>
                            </div>
                            <div>
                                <p><strong>Hash:</strong> <span style={{ wordBreak: 'break-all' }}>{block.hash}</span></p>
                                <p><strong>Previous Hash:</strong> <span style={{ wordBreak: 'break-all' }}>{block.previousHash}</span></p>
                            </div>
                        </div>
                    </div>

                    <div style={{ marginTop: '15px' }}>
                        <h4 style={{ margin: '0 0 10px 0', color: '#2c3e50' }}>Transactions:</h4>
                        {block.transactions.length === 0 ? (
                            <p style={{ color: '#666', fontStyle: 'italic' }}>No transactions in this block</p>
                        ) : (
                            <div style={{
                                display: 'grid',
                                gap: '10px'
                            }}>
                                {block.transactions.map((tx, idx) => (
                                    <div
                                        key={idx}
                                        style={{
                                            border: '1px solid #e9ecef',
                                            padding: '10px',
                                            borderRadius: '4px',
                                            backgroundColor: '#fff'
                                        }}
                                    >
                                        <p><strong>From:</strong> {tx.from}</p>
                                        <p><strong>To:</strong> {tx.to}</p>
                                        <p><strong>Amount:</strong> {tx.amount}</p>
                                        <p><strong>Message:</strong> {tx.message}</p>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            ))}
        </div>
    );
}

export default Blockchain;