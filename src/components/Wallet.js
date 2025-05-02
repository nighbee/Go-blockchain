import React, { useState } from 'react';
import { getWalletBalance } from './api';

function Wallet() {
    const [address, setAddress] = useState('');
    const [balance, setBalance] = useState(null);

    const handleCheckBalance = async () => {
        const data = await getWalletBalance(address);
        setBalance(data.balance);
    };

    return (
        <div>
            <h2>Wallet</h2>
            <input
                type="text"
                placeholder="Enter Wallet Address"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
            />
            <button onClick={handleCheckBalance}>Check Balance</button>
            {balance !== null && <p>Balance: {balance}</p>}
        </div>
    );
}

export default Wallet;