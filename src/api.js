export const getBlockchain = async () => {
    try {
        const response = await fetch('http://localhost:5001/chain', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status} :: GET /chain`);
        }
        const data = await response.json();
        console.log('Blockchain data:', data);
        return data;
    } catch (error) {
        console.error('Fetch blockchain error:', error);
        throw error;
    }
};

export const createWallet = async (blockchainAddress = '') => {
    try {
        const response = await fetch('http://localhost:5001/wallet/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(blockchainAddress ? { blockchainAddress } : {}),
        });
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(`HTTP error! status: ${response.status} :: ${errorData.message || 'Wallet creation failed'}`);
        }
        const data = await response.json();
        console.log('Wallet created:', data);
        return data;
    } catch (error) {
        console.error('Create wallet error:', error);
        throw error;
    }
};

export const createTransaction = async (transactionData) => {
    try {
        console.log('Transaction request:', transactionData);
        const response = await fetch('http://localhost:5001/transactions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(transactionData),
        });
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(`HTTP error! status: ${response.status} :: ${errorData.message || 'Transaction creation failed'}`);
        }
        const data = await response.json();
        console.log('Transaction created:', data);
        return { ...transactionData, status: data.status || 'success' };
    } catch (error) {
        console.error('Create transaction error:', error);
        throw error;
    }
};

export const getNodes = async () => {
    try {
        const response = await fetch('http://localhost:5001/nodes', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status} :: GET /nodes`);
        }
        const data = await response.json();
        console.log('Nodes data:', data);
        return data;
    } catch (error) {
        console.error('Fetch nodes error:', error);
        throw error;
    }
}; 