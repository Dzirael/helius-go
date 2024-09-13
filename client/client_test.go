package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/test-go/testify/assert"
)

const apiKey = "cd0664da-339d-48e5-a22e-f7c214d81f6c"

func TestParsedTransaction(t *testing.T) {
	service := New(apiKey)
	tx, err := service.GetParsedTransaction(context.Background(), "4yNVSuFLjxVaNxk4RPSz8aL67jNQshLZLiMh3tcaLR1G8DnSEnvAGrG8HP35gRfeiESRWmqWShGg3AKLVVe5QUre")
	assert.NoError(t, err)
	assert.NotEmpty(t, tx)
}

func TestGetTransactionHistory(t *testing.T) {
	service := New(apiKey)
	tx, err := service.GetTransactionHistory(context.Background(), nil, "7DtLV4eguy7Eg9Wf6AehfjfjSZNFwYE4TGqk4Fob6vHE")
	assert.NoError(t, err)
	assert.NotEmpty(t, tx)
}

func TestGetAllTransactionHistory(t *testing.T) {
	service := New(apiKey)
	txs := service.GetAllTransactionHistory(context.Background(), nil, "7DtLV4eguy7Eg9Wf6AehfjfjSZNFwYE4TGqk4Fob6vHE")
	assert.NotEmpty(t, txs)

	for txs.Next() {
		fmt.Println(txs.Result())
	}
}
