import React, { useState, useContext } from 'react';
import { WalletContext } from './WalletContext';
import { TransactionContext } from './TransactionContext';

function Transactions() {
    const { wallets } = useContext(WalletContext);
    const { transactions, addTransaction } = useContext(TransactionContext);
    const [senderAddress, setSenderAddress] = useState('');
    const [recipientAddress, setRecipientAddress] = useState('');
    const [message, setMessage] = useState('');
    const [amount, setAmount] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleSubmitTransaction = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            // Validate inputs
            if (!senderAddress || !recipientAddress || !message || !amount || amount <= 0) {
                throw new Error('Please fill all fields with valid values');
            }

            if (senderAddress === recipientAddress) {
                throw new Error('Sender and recipient cannot be the same');
            }

            // Get sender wallet
            const senderWallet = wallets.find(w => w.address === senderAddress);
            if (!senderWallet || !senderWallet.public_key || !senderWallet.private_key) {
                throw new Error('Invalid sender wallet or missing keys');
            }

            // Get recipient wallet
            const recipientWallet = wallets.find(w => w.address === recipientAddress);
            if (!recipientWallet) {
                throw new Error('Invalid recipient wallet');
            }

            // Create transaction object for signing
            const transactionData = {
                senderBlockchainAddress: senderWallet.address,
                recipientBlockchainAddress: recipientWallet.address,
                message: message,
                value: parseFloat(amount),
                privateKey: senderWallet.private_key,
                publicKey: senderWallet.public_key
            };

            // Step 1: Get signature from server
            const signResponse = await fetch('http://localhost:5001/sign', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(transactionData)
            });

            if (!signResponse.ok) {
                const errorData = await signResponse.json();
                throw new Error(errorData.message || 'Failed to sign transaction');
            }

            const { signature } = await signResponse.json();

            // Step 2: Submit transaction with signature
            const transactionPayload = {
                senderBlockchainAddress: senderWallet.address,
                recipientBlockchainAddress: recipientWallet.address,
                message: message,
                value: parseFloat(amount),
                senderPublicKey: senderWallet.public_key,
                signature: signature
            };

            const submitResponse = await fetch('http://localhost:5001/transactions', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(transactionPayload)
            });

            if (!submitResponse.ok) {
                const errorData = await submitResponse.json();
                throw new Error(errorData.message || 'Failed to submit transaction');
            }

            const newTransaction = await submitResponse.json();
            
            // Add transaction to global state
            addTransaction({
                ...transactionPayload,
                status: 'success'
            });

            // Clear form
            setMessage('');
            setAmount('');
            setSenderAddress('');
            setRecipientAddress('');

        } catch (err) {
            setError(err.message);
            console.error('Transaction error:', err);
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
                            placeholder="Enter transaction message"
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
                            placeholder="Enter amount"
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
                    {loading ? 'Processing...' : 'Send Transaction'}
                </button>
            </form>

            {error && (
                <div style={{ color: 'red', marginBottom: '20px' }}>
                    {error}
                </div>
            )}

            <h3>Transaction History</h3>
            {transactions.length === 0 ? (
                <p>No transactions yet</p>
            ) : (
                <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
                    {transactions.map((tx, index) => (
                        <div
                            key={index}
                            style={{
                                border: '1px solid #ddd',
                                padding: '15px',
                                marginBottom: '10px',
                                borderRadius: '5px',
                                backgroundColor: '#f9f9f9'
                            }}
                        >
                            <p><strong>Time:</strong> {tx.timestamp}</p>
                            <p><strong>From:</strong> {tx.senderBlockchainAddress}</p>
                            <p><strong>To:</strong> {tx.recipientBlockchainAddress}</p>
                            <p><strong>Amount:</strong> {tx.value}</p>
                            <p><strong>Message:</strong> {tx.message}</p>
                            <p><strong>Status:</strong> {tx.status || 'Pending'}</p>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

export default Transactions;