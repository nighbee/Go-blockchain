import React, { createContext, useState, useEffect } from 'react';

export const TransactionContext = createContext();

export const TransactionProvider = ({ children }) => {
    const [transactions, setTransactions] = useState(() => {
        const savedTransactions = localStorage.getItem('blockchain_transactions');
        return savedTransactions ? JSON.parse(savedTransactions) : [];
    });

    const [miningHistory, setMiningHistory] = useState(() => {
        const savedMiningHistory = localStorage.getItem('blockchain_mining_history');
        return savedMiningHistory ? JSON.parse(savedMiningHistory) : [];
    });

    // Save transactions to localStorage whenever they change
    useEffect(() => {
        localStorage.setItem('blockchain_transactions', JSON.stringify(transactions));
    }, [transactions]);

    // Save mining history to localStorage whenever it changes
    useEffect(() => {
        localStorage.setItem('blockchain_mining_history', JSON.stringify(miningHistory));
    }, [miningHistory]);

    const addTransaction = (transaction) => {
        setTransactions(prev => [{
            ...transaction,
            timestamp: new Date().toLocaleString(),
            status: 'success'
        }, ...prev]);
    };

    const addMiningRecord = (record) => {
        setMiningHistory(prev => [{
            ...record,
            timestamp: new Date().toLocaleString()
        }, ...prev]);
    };

    return (
        <TransactionContext.Provider value={{ 
            transactions, 
            miningHistory,
            addTransaction,
            addMiningRecord
        }}>
            {children}
        </TransactionContext.Provider>
    );
}; 