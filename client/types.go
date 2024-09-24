package client

import (
	"reflect"
	"strings"
)

type TransactionType string
type TransactionSource string
type Commitment string

type TransactionQuerry struct {
	Type       TransactionType
	Source     TransactionSource
	Before     string
	After      string
	Commitment Commitment
	Limit      string
}

func (q TransactionQuerry) ToMap() map[string]string {
	result := make(map[string]string)
	value := reflect.ValueOf(q)
	typ := reflect.TypeOf(q)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typ.Field(i).Name

		queryParam := fieldName
		if field.Kind() == reflect.String && field.String() != "" {
			result[strings.ToLower(queryParam)] = field.String()
		}
	}

	return result
}

type TransactionHistory struct {
	channel chan Transaction
	err     error
	current Transaction
}

func (t *TransactionHistory) Next() bool {
	transaction, ok := <-t.channel
	if ok {
		t.current = transaction
		return true
	}
	return false
}

func (t *TransactionHistory) Result() Transaction {
	return t.current
}

func (t *TransactionHistory) Err() error {
	return t.err
}

type Transaction struct {
	Description      string            `json:"description"`
	Type             string            `json:"type"`
	Source           string            `json:"source"`
	Fee              int               `json:"fee"`
	FeePayer         string            `json:"feePayer"`
	Signature        string            `json:"signature"`
	Slot             int               `json:"slot"`
	Timestamp        int               `json:"timestamp"`
	NativeTransfers  []NativeTransfer  `json:"nativeTransfers"`
	TokenTransfers   []TokenTransfer   `json:"tokenTransfers"`
	AccountData      []AccountData     `json:"accountData"`
	TransactionError *TransactionError `json:"transactionError"`
	Instructions     []Instruction     `json:"instructions"`
	Events           Events            `json:"events"`
}

type NativeTransfer struct {
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
}

type TokenTransfer struct {
	FromUserAccount  string  `json:"fromUserAccount"`
	ToUserAccount    string  `json:"toUserAccount"`
	FromTokenAccount string  `json:"fromTokenAccount"`
	ToTokenAccount   string  `json:"toTokenAccount"`
	TokenAmount      float64 `json:"tokenAmount"`
	Mint             string  `json:"mint"`
}

type AccountData struct {
	Account             string               `json:"account"`
	NativeBalanceChange int                  `json:"nativeBalanceChange"`
	TokenBalanceChanges []TokenBalanceChange `json:"tokenBalanceChanges"`
}

type TokenBalanceChange struct {
	UserAccount    string      `json:"userAccount"`
	TokenAccount   string      `json:"tokenAccount"`
	Mint           string      `json:"mint"`
	RawTokenAmount TokenAmount `json:"rawTokenAmount"`
}

type TokenAmount struct {
	TokenAmount string `json:"tokenAmount"`
}

type TransactionError struct {
	Error string `json:"error"`
}

type Instruction struct {
	Accounts          []string           `json:"accounts"`
	Data              string             `json:"data"`
	ProgramID         string             `json:"programId"`
	InnerInstructions []InnerInstruction `json:"innerInstructions"`
}

type InnerInstruction struct {
	Accounts  []string `json:"accounts"`
	Data      string   `json:"data"`
	ProgramID string   `json:"programId"`
}

type Events struct {
	Nft  NftEvent  `json:"nft"`
	Swap SwapEvent `json:"swap"`
}

type NftEvent struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	Amount      int    `json:"amount"`
	Fee         int    `json:"fee"`
	FeePayer    string `json:"feePayer"`
	Signature   string `json:"signature"`
	Slot        int    `json:"slot"`
	Timestamp   int    `json:"timestamp"`
	SaleType    string `json:"saleType"`
	Buyer       string `json:"buyer"`
	Seller      string `json:"seller"`
	Staker      string `json:"staker"`
	Nfts        []Nft  `json:"nfts"`
}

type Nft struct {
	Mint          string `json:"mint"`
	TokenStandard string `json:"tokenStandard"`
}

type SwapEvent struct {
	NativeInput  NativeAmount  `json:"nativeInput"`
	NativeOutput NativeAmount  `json:"nativeOutput"`
	TokenInputs  []TokenInput  `json:"tokenInputs"`
	TokenOutputs []TokenOutput `json:"tokenOutputs"`
	TokenFees    []TokenFee    `json:"tokenFees"`
	NativeFees   []NativeFee   `json:"nativeFees"`
	InnerSwaps   []InnerSwap   `json:"innerSwaps"`
}

type NativeAmount struct {
	Account string `json:"account"`
	Amount  string `json:"amount"`
}

type TokenInput struct {
	UserAccount    string      `json:"userAccount"`
	TokenAccount   string      `json:"tokenAccount"`
	Mint           string      `json:"mint"`
	RawTokenAmount TokenAmount `json:"rawTokenAmount"`
}

type TokenOutput struct {
	UserAccount    string      `json:"userAccount"`
	TokenAccount   string      `json:"tokenAccount"`
	Mint           string      `json:"mint"`
	RawTokenAmount TokenAmount `json:"rawTokenAmount"`
}

type TokenFee struct {
	UserAccount    string      `json:"userAccount"`
	TokenAccount   string      `json:"tokenAccount"`
	Mint           string      `json:"mint"`
	RawTokenAmount TokenAmount `json:"rawTokenAmount"`
}

type NativeFee struct {
	Account string `json:"account"`
	Amount  string `json:"amount"`
}

type InnerSwap struct {
	TokenInputs  []TokenInput  `json:"tokenInputs"`
	TokenOutputs []TokenOutput `json:"tokenOutputs"`
	TokenFees    []TokenFee    `json:"tokenFees"`
	NativeFees   []NativeFee   `json:"nativeFees"`
	ProgramInfo  ProgramInfo   `json:"programInfo"`
}

type ProgramInfo struct {
	Source          string `json:"source"`
	Account         string `json:"account"`
	ProgramName     string `json:"programName"`
	InstructionName string `json:"instructionName"`
}
