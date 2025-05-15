// src/components/Transactions.js
import React, { useState, useContext } from 'react';
import { createTransaction } from './api';
import { WalletContext } from './WalletContext';

function Transactions() {
    const { wallets } = useContext(WalletContext);
    const [transactions, setTransactions] = useState([]);
    const [senderAddress, setSenderAddress] = useState('');
    const [recipientAddress, setRecipientAddress] = useState('');
    const [message, setMessage] = useState('');
    const [amount, setAmount] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleSubmitTransaction = async (e) => {
        e.preventDefault();
        if (!senderAddress || !recipientAddress || !message || !amount || amount <= 0) {
            setError('Please fill all fields with valid values');
            return;
        }
        if (senderAddress === recipientAddress) {
            setError('Sender and recipient cannot be the same');
            return;
        }
        setLoading(true);
        setError(null);
        try {
            const transaction = await createTransaction({
                senderBlockchainAddress: senderAddress,
                recipientBlockchainAddress: recipientAddress,
                message,
                value: parseFloat(amount),
                senderPublicKey: 'dummy',
                signature: 'dummy',
            });
            setTransactions([...transactions, transaction]);
            setMessage('');
            setAmount('');
            setSenderAddress('');
            setRecipientAddress('');
        } catch (err) {
            setError('Failed to create transaction');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ margin: '20px' }}>
            <h2>Create Transaction</h2>
            <form onSubmit={handleSubmitTransaction} style={{ marginBottom: '20px' }}>
                <div style={{ marginBottom: '10px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>
                        Sender Address:
                        <select
                            value={senderAddress}
                            onChange={(e) => setSenderAddress(e.target.value)}
                            style={{ width: '100%', padding: '8px', marginTop: '5px' }}
                        >
                            <option value="">Select a wallet</option>
                            {wallets.map((wallet, index) => (
                                <option key={index} value={wallet.address}>
                                    Wallet {index + 1}: {wallet.address}
                                </option>
                            ))}
                        </select>
                    </label>
                </div>
                <div style={{ marginBottom: '10px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>
                        Recipient Address:
                        <select
                            value={recipientAddress}
                            onChange={(e) => setRecipientAddress(e.target.value)}
                            style={{ width: '100%', padding: '8px', marginTop: '5px' }}
                        >
                            <option value="">Select a wallet</option>
                            {wallets.map((wallet, index) => (
                                <option key={index} value={wallet.address}>
                                    Wallet {index + 1}: {wallet.address}
                                </option>
                            ))}
                        </select>
                    </label>
                </div>
                <div style={{ marginBottom: '10px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>
                        Message:
                        <input
                            type="text"
                            value={message}
                            onChange={(e) => setMessage(e.target.value)}
                            placeholder="e.g., Payment for services"
                            style={{ width: '100%', padding: '8px', marginTop: '5px' }}
                        />
                    </label>
                </div>
                <div style={{ marginBottom: '10px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>
                        Amount:
                        <input
                            type="number"
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            placeholder="e.g., 0.5"
                            step="0.1"
                            min="0"
                            style={{ width: '100%', padding: '8px', marginTop: '5px' }}
                        />
                    </label>
                </div>
                <button
                    type="submit"
                    disabled={loading}
                    style={{
                        padding: '10px 20px',
                        backgroundColor: '#007bff',
                        color: 'white',
                        border: 'none',
                        borderRadius: '5px',
                        cursor: loading ? 'not-allowed' : 'pointer',
                    }}
                >
                    {loading ? 'Submitting...' : 'Send Transaction'}
                </button>
            </form>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            <h3>Submitted Transactions</h3>
            {transactions.length === 0 ? (
                <p>No transactions submitted yet.</p>
            ) : (
                <ul style={{ listStyle: 'none', padding: 0, margin: '0px', width: '100%' }}>
                    {transactions.map((tx, index) => (
                        <li
                            key={index}
                            style={{
                                border: '1px solid #ccc',
                                padding: '10px',
                                margin: '5px 0',
                                borderRadius: '5px',
                            }}
                        >
                            <p><strong>From:</strong> {tx.senderBlockchainAddress}</p>
                            <p><strong>To:</strong> {tx.recipientBlockchainAddress}</p>
                            <p><strong>Amount:</strong> {tx.value}</p>
                            <p><strong>Message:</strong> {tx.message}</p>
                            <p><strong>Status:</strong> {tx.status}</p>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Transactions;