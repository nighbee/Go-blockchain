import React, { useEffect, useState } from 'react';
import { getNodes } from '../api';

function Nodes() {
    const [nodes, setNodes] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchNodes = async () => {
            try {
                const data = await getNodes();
                setNodes(data.nodes || []);
            } catch (err) {
                console.error('Error fetching nodes:', err);
                setError('Failed to load nodes');
            } finally {
                setLoading(false);
            }
        };
        fetchNodes();
    }, []);

    if (loading) {
        return (
            <div style={{ margin: '20px' }}>
                <h2>Nodes</h2>
                <p>Loading nodes...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div style={{ margin: '20px' }}>
                <h2>Nodes</h2>
                <p style={{ color: 'red' }}>{error}</p>
            </div>
        );
    }

    return (
        <div style={{ margin: '20px' }}>
            <h2>Nodes</h2>
            {nodes.length === 0 ? (
                <p>No nodes connected</p>
            ) : (
                <ul style={{ listStyle: 'none', padding: 0 }}>
                    {nodes.map((node, index) => (
                        <li
                            key={index}
                            style={{
                                border: '1px solid #ddd',
                                padding: '15px',
                                margin: '10px 0',
                                borderRadius: '5px',
                                backgroundColor: '#f8f9fa'
                            }}
                        >
                            <p><strong>Node {index + 1}:</strong> {node}</p>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Nodes;