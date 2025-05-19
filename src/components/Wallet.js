import React, { useState, useContext, useEffect } from 'react';
import { WalletContext } from './WalletContext';
import { createWallet } from '../api';

const Wallet = () => {
    const { wallets, registerWallet } = useContext(WalletContext);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

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
                    data.wallets.forEach(address => {
                        // Create a wallet object with the address
                        const wallet = {
                            address: address,
                            publicKey: '', // These will be empty for existing wallets
                            privateKey: '' // These will be empty for existing wallets
                        };
                        registerWallet(wallet);
                    });
                }
            } catch (err) {
                console.error('Error fetching wallets:', err);
                setError('Failed to fetch existing wallets');
            }
        };

        fetchExistingWallets();
    }, [registerWallet]);

    const handleCreateWallet = async () => {
        setLoading(true);
        setError(null);
        try {
            const wallet = await createWallet();
            registerWallet(wallet);
        } catch (err) {
            console.error('Error creating wallet:', err);
            setError('Failed to create wallet');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold mb-8">Wallet Management</h1>
            
            <div className="mb-8">
                <button
                    onClick={handleCreateWallet}
                    disabled={loading}
                    className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
                >
                    {loading ? 'Creating...' : 'Create New Wallet'}
                </button>
                {error && <p className="text-red-500 mt-2">{error}</p>}
            </div>

            <div className="grid gap-6">
                {wallets.map((wallet, index) => (
                    <div key={index} className="bg-white p-6 rounded-lg shadow-md">
                        <h2 className="text-xl font-semibold mb-4">Wallet {index + 1}</h2>
                        <div className="space-y-2">
                            <p><span className="font-medium">Address:</span> {wallet.address}</p>
                            {wallet.publicKey && (
                                <p><span className="font-medium">Public Key:</span> {wallet.publicKey}</p>
                            )}
                            {wallet.privateKey && (
                                <p><span className="font-medium">Private Key:</span> {wallet.privateKey}</p>
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Wallet;