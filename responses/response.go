package responses

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
