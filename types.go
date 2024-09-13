package helius

import "github.com/Dzirael/helius-go-sdk/client"

type Deposits struct {
	channel chan client.Transaction
	err     error
	current client.Transaction
}

func (d *Deposits) Next() bool {
	transaction, ok := <-d.channel
	if ok {
		d.current = transaction
		return true
	}
	return false
}

func (d *Deposits) Result() client.Transaction {
	return d.current
}

func (d *Deposits) Err() error {
	return d.err
}
