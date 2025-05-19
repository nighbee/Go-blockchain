import React, { createContext, useState, useEffect } from 'react';
import { createWallet } from '../api';

export const WalletContext = createContext();

export const WalletProvider = ({ children }) => {
    const [wallets, setWallets] = useState([]);

    // Fetch existing wallets from the blockchain when component mounts
    useEffect(() => {
        const fetchExistingWallets = async () => {
            try {
                const response = await fetch('http://localhost:5001/wallets');
                if (!response.ok) {
                    throw new Error('Failed to fetch wallets');
                }
                const data = await response.json();
                // Add existing wallets to the context
                if (data.wallets && Array.isArray(data.wallets)) {
                    const walletObjects = data.wallets.map(address => ({
                        address: address,
                        publicKey: '', // These will be empty for existing wallets
                        privateKey: '' // These will be empty for existing wallets
                    }));
                    setWallets(walletObjects);
                }
            } catch (err) {
                console.error('Error fetching wallets:', err);
            }
        };

        fetchExistingWallets();
    }, []);

    const registerWallet = async (wallet) => {
        try {
            // If the wallet is just an address (from existing wallets), use it as is
            if (typeof wallet === 'string') {
                wallet = {
                    address: wallet,
                    publicKey: '',
                    privateKey: ''
                };
            }
            
            // Check if wallet already exists
            const exists = wallets.some(w => w.address === wallet.address);
            if (!exists) {
                setWallets(prevWallets => [...prevWallets, wallet]);
            }
            return wallet;
        } catch (err) {
            console.error('Failed to register wallet:', err);
            throw err;
        }
    };

    return (
        <WalletContext.Provider value={{ wallets, registerWallet }}>
            {children}
        </WalletContext.Provider>
    );
};