package models

type Account struct {
	AccountNumber string  `bson:"account_number" json:"account_number"`
	Name          string  `bson:"name" json:"name"`
	Pin           string  `bson:"pin" json:"-"`
	Balance       float64 `bson:"balance" json:"balance"`
}

type Transaction struct {
	From     string  `bson:"from" json:"from"`
	To       string  `bson:"to" json:"to"`
	Type     string  `bson:"type" json:"type"`
	Amount   float64 `bson:"amount" json:"amount"`
	DateTime string  `bson:"datetime" json:"datetime"`
}

type CreateAccountRequest struct {
	Name string `json:"name" binding:"required"`
	Pin  string `json:"pin" binding:"required"`
}

type DepositRequest struct {
	AccountNumber string  `json:"account_number" binding:required`
	Pin           string  `json:"pin" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

type WithdrawRequest struct {
	AccountNumber string  `json:"account_number" binding:"required"`
	Pin           string  `json:"pin" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

type TransferRequest struct {
	FromAccount string  `json:"from_account" binding:"required"`
	FromPin     string  `json:"from_pin" binding:"required"`
	ToAccount   string  `json:"to_account" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

type PinRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	OldPin        string `json:"old_pin" binding:"required"`
	NewPin        string `json:"new_pin" binding:"required"`
}
