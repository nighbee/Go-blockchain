import React, { createContext, useState, useEffect } from 'react';
import { createWallet } from '../api';

export const WalletContext = createContext();

export const WalletProvider = ({ children }) => {
    const [wallets, setWallets] = useState(() => {
        try {
            const savedWallets = localStorage.getItem('blockchain_wallets');
            return savedWallets ? JSON.parse(savedWallets) : [];
        } catch (err) {
            console.error('Error loading wallets from localStorage:', err);
            return [];
        }
    });

    // Save wallets to localStorage whenever they change
    useEffect(() => {
        try {
            localStorage.setItem('blockchain_wallets', JSON.stringify(wallets));
        } catch (err) {
            console.error('Error saving wallets to localStorage:', err);
        }
    }, [wallets]);

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