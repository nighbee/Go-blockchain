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
                console.log('Fetched blockchain data:', data); // Debugging log
                if (data.chain && data.chain.length > 0) {
                    setBlocks(data.chain); // Ensure the genesis block is included
                } else {
                    setError('Blockchain data is empty');
                }
            } catch (err) {
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
        <div>
            <h2>Blockchain Activity</h2>
            {blocks.map((block, index) => (
                <div key={index} style={{ border: '1px solid black', margin: '10px', padding: '10px' }}>
                    <p><strong>Block Index:</strong> {block.index}</p>
                    <p><strong>Timestamp:</strong> {new Date(block.timestamp).toLocaleString()}</p>
                    <p><strong>Hash:</strong> {block.hash}</p>
                    <p><strong>Transactions:</strong></p>
                    <ul>
                        {(block.transactions || []).map((tx, idx) => (
                            <li key={idx}>
                                <p><strong>From:</strong> {tx.from}</p>
                                <p><strong>To:</strong> {tx.to}</p>
                                <p><strong>Amount:</strong> {tx.amount}</p>
                                <p><strong>Message:</strong> {tx.message}</p>
                            </li>
                        ))}
                    </ul>
                </div>
            ))}
        </div>
    );
}

export default Blockchain;