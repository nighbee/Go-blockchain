import React, { useEffect, useState } from 'react';
import { getNodes } from './api';

function Nodes() {
    const [nodes, setNodes] = useState([]);

    useEffect(() => {
        const fetchNodes = async () => {
            const data = await getNodes();
            setNodes(data.nodes);
        };
        fetchNodes();
    }, []);

    return (
        <div>
            <h2>Nodes</h2>
            <ul>
                {nodes.map((node, index) => (
                    <li key={index}>{node}</li>
                ))}
            </ul>
        </div>
    );
}

export default Nodes;