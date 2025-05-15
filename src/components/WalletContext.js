// src/WalletContext.js
import React, { createContext, useState } from 'react';

export const WalletContext = createContext();

export const WalletProvider = ({ children }) => {
    const [wallets, setWallets] = useState([]);

    return (
        <WalletContext.Provider value={{ wallets, setWallets }}>
            {children}
        </WalletContext.Provider>
    );
};