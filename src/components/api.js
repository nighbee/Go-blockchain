// src/api.js
export const getBlockchain = async () => {
    try {
        const response = await fetch('http://localhost:5001/chain', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log('Raw response:', data);
        return data;
    } catch (error) {
        console.error('Fetch error:', error);
        throw error;
    }
};

export const createWallet = async () => {
    try {
        const address = `wallet-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
        const response = await fetch('http://localhost:5001/wallet/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ blockchainAddress: address }),
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log('Wallet created:', data);
        return { address, message: data.message };
    } catch (error) {
        console.error('Create wallet error:', error);
        throw error;
    }
};

export const createTransaction = async (transactionData) => {
    try {
        const response = await fetch('http://localhost:5001/transactions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(transactionData),
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log('Transaction created:', data);
        return { ...transactionData, status: data.status || 'success' };
    } catch (error) {
        console.error('Create transaction error:', error);
        throw error;
    }
};