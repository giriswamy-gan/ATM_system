package main

import (
	"context"
<<<<<<< HEAD
	"log"
=======
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
>>>>>>> upload
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Account struct {
	AccountNumber string  `bson:"account_number" json:"account_number"`
<<<<<<< HEAD
=======
	Name          string  `bson:"name" json:"name"`
>>>>>>> upload
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

var client *mongo.Client

<<<<<<< HEAD
=======
func hashPassword(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}

>>>>>>> upload
func main() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Ping the primary
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Gin router
	r := gin.Default()

	// Create an account
	type CreateAccountRequest struct {
<<<<<<< HEAD
		AccountNumber string `json:"account_number" binding:"required"`
		Pin           string `json:"pin" binding:"required"`
=======
		Name string `json:"name" binding:"required"`
		Pin  string `json:"pin" binding:"required"`
>>>>>>> upload
	}
	r.POST("/create", func(c *gin.Context) {
		var req CreateAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

<<<<<<< HEAD
		// Check if account exists
		filter := bson.M{"account_number": req.AccountNumber}
		var existingAccount Account
		err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&existingAccount)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "account already exists"})
=======
		accountNumber := fmt.Sprintf("%06d", rand.Intn(1000000))

		// Check if account exists
		// filter := bson.M{"account_number": accountNumber}
		// var existingAccount Account
		// err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&existingAccount)
		// if err == nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "account already exists"})
		// 	return
		// }

		if len(req.Pin) != 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pin must be 4 digits"})
>>>>>>> upload
			return
		}

		// Insert new account
		newAccount := Account{
<<<<<<< HEAD
			AccountNumber: req.AccountNumber,
			Pin:           req.Pin,
=======
			Name:          req.Name,
			AccountNumber: accountNumber,
			Pin:           hashPassword(req.Pin),
>>>>>>> upload
			Balance:       0,
		}
		_, err = client.Database("atm").Collection("accounts").InsertOne(context.Background(), newAccount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "account created successfully"})
	})

	// Deposit money into an account
	type DepositRequest struct {
		AccountNumber string  `json:"account_number" binding:required`
		Pin           string  `json:"pin" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}
	r.POST("/deposit", func(c *gin.Context) {
		var req DepositRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if account exists and PIN matches
<<<<<<< HEAD
		filter := bson.M{"account_number": req.AccountNumber, "pin": req.Pin}
=======
		filter := bson.M{"account_number": req.AccountNumber, "pin": hashPassword(req.Pin)}
>>>>>>> upload
		var existingAccount Account
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
		transaction := Transaction{
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
	})

	//withdraw money
	type WithdrawRequest struct {
		AccountNumber string  `json:"account_number" binding:"required"`
		Pin           string  `json:"pin" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}
	r.POST("/withdraw", func(c *gin.Context) {
		var req WithdrawRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//check if account exists and pin match
<<<<<<< HEAD
		filter := bson.M{"account_number": req.AccountNumber, "pin": req.Pin}
=======
		filter := bson.M{"account_number": req.AccountNumber, "pin": hashPassword(req.Pin)}
>>>>>>> upload
		var existingAccount Account
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
		transaction := Transaction{
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

	})

	// Transfer money from one account to another
	type TransferRequest struct {
		FromAccount string  `json:"from_account" binding:"required"`
		FromPin     string  `json:"from_pin" binding:"required"`
		ToAccount   string  `json:"to_account" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
	}
	r.POST("/transfer", func(c *gin.Context) {
		var req TransferRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if 'from' account exists and PIN matches
<<<<<<< HEAD
		filterFrom := bson.M{"account_number": req.FromAccount, "pin": req.FromPin}
=======
		filterFrom := bson.M{"account_number": req.FromAccount, "pin": hashPassword(req.FromPin)}
>>>>>>> upload
		var fromAccount Account
		err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filterFrom).Decode(&fromAccount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' account number or PIN"})
			return
		}

		// Check if 'to' account exists
		filterTo := bson.M{"account_number": req.ToAccount}
		var toAccount Account
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
		transactionFrom := Transaction{
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

		transactionTo := Transaction{
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
	})

	// Set or reset PIN
	type PinRequest struct {
		AccountNumber string `json:"account_number" binding:"required"`
		OldPin        string `json:"old_pin" binding:"required"`
		NewPin        string `json:"new_pin" binding:"required"`
	}
	r.POST("/setpin", func(c *gin.Context) {
		var req PinRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if account exists and old PIN matches
<<<<<<< HEAD
		filter := bson.M{"account_number": req.AccountNumber, "pin": req.OldPin}
=======
		filter := bson.M{"account_number": req.AccountNumber, "pin": hashPassword(req.OldPin)}
>>>>>>> upload
		var account Account
		err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&account)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account number or PIN"})
			return
		}

		// Update PIN
<<<<<<< HEAD
		update := bson.M{"$set": bson.M{"pin": req.NewPin}}
=======
		update := bson.M{"$set": bson.M{"pin": hashPassword(req.NewPin)}}
>>>>>>> upload
		_, err = client.Database("atm").Collection("accounts").UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update PIN"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "PIN updated successfully"})
	})

	//get transaction history of account
	r.POST("/bankstatement", func(c *gin.Context) {
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
		var account Account
		err := client.Database("atm").Collection("accounts").FindOne(context.Background(), filter).Decode(&account)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account number"})
			return
		}

		// Validate PIN
<<<<<<< HEAD
		if account.Pin != req.Pin {
=======
		if account.Pin != hashPassword(req.Pin) {
>>>>>>> upload
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
		var transactions []Transaction
		for cursor.Next(context.Background()) {
			var transaction Transaction
			if err := cursor.Decode(&transaction); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode transaction"})
				return
			}
			transactions = append(transactions, transaction)
		}

		c.JSON(http.StatusOK, gin.H{"transactions": transactions})
	})

	r.Run("localhost:9003")
}
