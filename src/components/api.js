const API_BASE_URL = 'http://localhost:5001';

export const getBlockchain = async () => {
    const response = await fetch('http://localhost:5001/chain');
    if (!response.ok) {
        throw new Error('Failed to fetch blockchain data');
    }
    return await response.json();
};
export const getNodes = async () => {
    const response = await fetch('http://localhost:5001/chain');
    if (!response.ok) {
        throw new Error('Failed to fetch nodes data');
    }
    return await response.json();
};

export const getWalletBalance = async (address) => {
    const response = await fetch(`${API_BASE_URL}/api/wallet/${address}`);
    if (!response.ok) {
        throw new Error('Failed to fetch wallet balance');
    }
    return await response.json();
};