package controllers

import (
	"atm-system/configs"
	"atm-system/models"
	"atm-system/responses"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var client = configs.DB

func hashPassword(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}

// creating an account

func CreateAccount(c *gin.Context) {
	var req responses.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountNumber := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Check if account exists
	filter := bson.M{"account_number": accountNumber}
	var existingAccount models.Account
	err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&existingAccount)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account already exists"})
		return
	}

	if len(req.Pin) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pin must be 4 digits"})
		return
	}

	// Insert new account
	newAccount := models.Account{
		Name:          req.Name,
		AccountNumber: accountNumber,
		Pin:           hashPassword(req.Pin),
		Balance:       0,
	}
	_, err = client.Database("atm").Collection("accounts").InsertOne(context.Background(), newAccount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Account created successfully. Account number": accountNumber})
}

// depositing money

func Deposit(c *gin.Context) {
	var req responses.DepositRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if account exists and PIN matches
	filter := bson.M{"account_number": req.AccountNumber, "pin": hashPassword(req.Pin)}
	var existingAccount models.Account
	err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&existingAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account number or PIN"})
		return
	}

	// Update account balance
	existingAccount.Balance += req.Amount
	update := bson.M{"$set": bson.M{"balance": existingAccount.Balance}}
	_, err = client.Database("atm").Collection("accounts").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deposit money"})
		return
	}

	// Insert transaction record
	transaction := models.Transaction{
		From:     "",
		To:       req.AccountNumber,
		Type:     "deposit",
		Amount:   req.Amount,
		DateTime: time.Now().Format(time.RFC3339),
	}
	_, err = client.Database("atm").Collection("transactions").InsertOne(context.Background(), transaction)
	if err != nil {
		log.Println("failed to insert transaction record:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "money deposited successfully"})
}

// withdraw money

func Withdraw(c *gin.Context) {
	var req responses.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if account exists and pin match
	filter := bson.M{"account_number": req.AccountNumber, "pin": hashPassword(req.Pin)}
	var existingAccount models.Account
	err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&existingAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account number or PIN"})
		return
	}

	// Check if 'from' account has enough balance
	if existingAccount.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough balance in account"})
		return
	}

	//update account balance
	existingAccount.Balance -= req.Amount
	update := bson.M{"$set": bson.M{"balance": existingAccount.Balance}}
	_, err = client.Database("atm").Collection("accounts").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to withdraw money"})
		return
	}

	//insert transaction record
	transaction := models.Transaction{
		From:     "",
		To:       req.AccountNumber,
		Type:     "withdraw",
		Amount:   req.Amount,
		DateTime: time.Now().Format(time.RFC3339),
	}
	_, err = client.Database("atm").Collection("transactions").InsertOne(context.Background(), transaction)
	if err != nil {
		log.Println("failed to insert transaction record:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "money withdraw successful"})

}

// transfer money from one account to another

func Transfer(c *gin.Context) {
	var req responses.TransferRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if 'from' account exists and PIN matches
	filterFrom := bson.M{"account_number": req.FromAccount, "pin": hashPassword(req.FromPin)}
	var fromAccount models.Account
	err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filterFrom).Decode(&fromAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' account number or PIN"})
		return
	}

	// Check if 'to' account exists
	filterTo := bson.M{"account_number": req.ToAccount}
	var toAccount models.Account
	err = client.Database("atm").Collection("accounts").FindOne(context.Background(), filterTo).Decode(&toAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' account number"})
		return
	}

	// Check if 'from' account has enough balance
	if fromAccount.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough balance in 'from' account"})
		return
	}

	// Update 'from' account balance
	fromAccount.Balance -= req.Amount
	updateFrom := bson.M{"$set": bson.M{"balance": fromAccount.Balance}}
	_, err = client.Database("atm").Collection("accounts").UpdateOne(context.Background(), filterFrom, updateFrom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to transfer money"})
		return
	}

	// Update 'to' account balance
	toAccount.Balance += req.Amount
	updateTo := bson.M{"$set": bson.M{"balance": toAccount.Balance}}
	_, err = client.Database("atm").Collection("accounts").UpdateOne(context.Background(), filterTo, updateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to transfer money"})
		return
	}

	// Insert transaction records
	transactionFrom := models.Transaction{
		From:     req.FromAccount,
		To:       req.ToAccount,
		Type:     "withdraw",
		Amount:   req.Amount,
		DateTime: time.Now().Format(time.RFC3339),
	}
	_, err = client.Database("atm").Collection("transactions").InsertOne(context.Background(), transactionFrom)
	if err != nil {
		log.Println("failed to insert transaction record:", err)
	}

	transactionTo := models.Transaction{
		From:     req.FromAccount,
		To:       req.ToAccount,
		Type:     "deposit",
		Amount:   req.Amount,
		DateTime: time.Now().Format(time.RFC3339),
	}
	_, err = client.Database("atm").Collection("transactions").InsertOne(context.Background(), transactionTo)
	if err != nil {
		log.Println("failed to insert transaction record:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "money transferred successfully"})
}

// changing pin

func SetPin(c *gin.Context) {
	var req responses.PinRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if account exists and old PIN matches
	filter := bson.M{"account_number": req.AccountNumber, "pin": hashPassword(req.OldPin)}
	var account models.Account
	err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account number or PIN"})
		return
	}

	// Update PIN
	update := bson.M{"$set": bson.M{"pin": hashPassword(req.NewPin)}}
	_, err = client.Database("atm").Collection("accounts").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN updated successfully"})
}

// show entire transaction history of an account

func BankStatement(c *gin.Context) {
	var req struct {
		AccountNumber string `json:"account_number"`
		Pin           string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Find account by account number
	filter := bson.M{"account_number": req.AccountNumber}
	var account models.Account
	err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account number"})
		return
	}

	// Validate PIN
	if account.Pin != hashPassword(req.Pin) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid PIN"})
		return
	}

	// Find all transactions for account
	filter = bson.M{"$or": []interface{}{
		bson.M{"from": req.AccountNumber},
		bson.M{"to": req.AccountNumber},
	}}
	cursor, err := client.Database("atm").Collection("transactions").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transaction history"})
		return
	}
	defer cursor.Close(context.Background())

	// Extract transactions from cursor
	var transactions []models.Transaction
	for cursor.Next(context.Background()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode transaction"})
			return
		}
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
