package helius

import (
	"context"

	"github.com/Dzirael/helius-go-sdk/client"
)

type Helius interface {
	// GetParsedTransaction returns a list of transactions
	//
	// Parameters:
	// - ctx: context
	// - transactions: list of transaction hashes
	//
	// Returns:
	// - []Transaction: list of transactions
	// - error: error
	GetParsedTransaction(ctx context.Context, transactions ...string) ([]client.Transaction, error)

	// GetTransactionHistory returns a list of transactions
	//
	// Parameters:
	// - ctx: context
	// - params: helius querry parameters
	// - address: address to check
	//
	// Returns:
	// - []Transaction: list of transactions
	// - error: error
	GetTransactionHistory(ctx context.Context, params *client.TransactionQuerry, address string) ([]client.Transaction, error)

	// GetAllTransactionHistory returns a channel of transaction can be used to iterate over all transactions
	//
	// Parameters:
	// - ctx: context
	// - params: helius querry parameters
	// - address: address to check
	//
	// Returns:
	// - TransactionHistory: channel of transactions
	//
	// Example:
	// 	transactions := helius.GetAllTransactionHistory(ctx, nil, "address")
	// 	for transactions.Next() {
	// 		tx := transactions.Result()
	// 		fmt.Println(tx)
	// 	}
	// 	if transactions.Err() != nil {
	// 		log.Fatal(transactions.Err())
	// 	}
	GetAllTransactionHistory(ctx context.Context, params *client.TransactionQuerry, address string) client.TransactionHistory
	// GetDeposits returns a channel of transactions that are deposits to the address
	// The channel will be closed when there are no more transactions
	// The error will be set if there was an error fetching the transactions
	//
	// Parameters:
	// - ctx: context
	// - address: address to check
	//
	// Returns:
	// - Deposits: channel of transactions
	//
	// Example:
	// 	deposits := helius.GetDeposits(ctx, "address")
	// 	for deposits.Next() {
	// 		tx := deposits.Result()
	// 		fmt.Println(tx)
	// 	}
	// 	if deposits.Err() != nil {
	// 		log.Fatal(deposits.Err())
	// 	}
	GetDeposits(ctx context.Context, address string, before *string) Deposits
}

type helius struct {
	httpClinet client.HttpClient
}

func New(apiKey string) Helius {
	return &helius{
		httpClinet: client.New(apiKey),
	}
}

func (h *helius) GetParsedTransaction(ctx context.Context, transactions ...string) ([]client.Transaction, error) {
	return h.httpClinet.GetParsedTransaction(ctx, transactions...)
}

func (h *helius) GetTransactionHistory(ctx context.Context, params *client.TransactionQuerry, address string) ([]client.Transaction, error) {
	return h.httpClinet.GetTransactionHistory(ctx, params, address)
}

func (h *helius) GetAllTransactionHistory(ctx context.Context, params *client.TransactionQuerry, address string) client.TransactionHistory {
	return h.httpClinet.GetAllTransactionHistory(ctx, params, address)
}

func (h *helius) GetDeposits(ctx context.Context, address string, before *string) Deposits {
	params := &client.TransactionQuerry{
		Type:   client.TRANSFER,
		Source: client.SYSTEM_PROGRAM,
	}
	if before != nil {
		params.Before = *before
	}

	result := Deposits{
		channel: make(chan client.Transaction),
	}

	go func() {
		defer close(result.channel)
		txs := h.httpClinet.GetAllTransactionHistory(ctx, params, address)
		for txs.Next() {
			tx := txs.Result()
			if tx.FeePayer != address && len(tx.NativeTransfers) == 1 {
				result.channel <- tx
			}
		}

		if txs.Err() != nil {
			result.err = txs.Err()
		}
	}()
	return result
}
