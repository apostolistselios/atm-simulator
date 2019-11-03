package api

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// User struct
type User struct {
	Balance int
	// Daily withdraw limit.
	WithdrawalLimit int
	// Amount withdrawed the current day.
	AmountWithdrawed int
	// The day the last withdraw transaction happened.
	LastWithdrawalDay int
}

// Transaction struct
type Transaction struct {
	UserID string
	Type   string
	Amount int
}

// GetDatabase returns a Database object and error.
func GetDatabase(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// ExecuteTransaction updates the balance of the user with that specific username.
// transaction slice [username, tranType, amount]
func ExecuteTransaction(db *bolt.DB, transaction Transaction) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(transaction.UserID))

		var u User
		json.Unmarshal(value, &u)
		if transaction.Type == "w" {
			if err := processWithdrawal(&u, &transaction); err != nil {
				return err
			}
		} else if transaction.Type == "d" {
			u.Balance += transaction.Amount
		}
		encodedUser, _ := json.Marshal(u)
		return bucket.Put([]byte(transaction.UserID), encodedUser)
	})
}

// processWithdrawal makes the required checks in order for the withdrawal to be carried on.
// If the checks pass the User is prepared and the withdrawal is ready to be executed.
// If an error occurs it returns the error, else it returns nil.
func processWithdrawal(u *User, transaction *Transaction) error {
	currDay := time.Now().Day()
	if u.Balance > 0 && u.Balance-transaction.Amount >= 0 {
		if currDay == u.LastWithdrawalDay && u.AmountWithdrawed+transaction.Amount > u.WithdrawalLimit {
			return fmt.Errorf("error withdrawal limit, Daily Limit: %d", u.WithdrawalLimit)
		} else if currDay == u.LastWithdrawalDay {
			u.Balance -= transaction.Amount
			u.AmountWithdrawed += transaction.Amount
		} else {
			u.Balance -= transaction.Amount
			u.AmountWithdrawed = transaction.Amount
			u.LastWithdrawalDay = currDay
		}
	} else {
		return fmt.Errorf("error high withdrawal amount, User: %s Balance: %d", transaction.UserID, u.Balance)
	}
	return nil
}

// GetUser returns the User with that specific username.
func GetUser(db *bolt.DB, username string) (User, error) {
	var u User
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(username))
		if value == nil {
			return fmt.Errorf("error User: %s doesn't exist", username)
		}

		json.Unmarshal(value, &u)
		return nil
	})

	return u, err
}

// GetBalance returns the balance of the specific user.
func GetBalance(db *bolt.DB, username string) (int, error) {
	var u User
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(username))

		json.Unmarshal(value, &u)
		return nil
	})
	return u.Balance, err
}

// ViewUsers prints all the users in the console.
func ViewUsers(db *bolt.DB) {
	fmt.Println("-------USERS-------")
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		bucket.ForEach(func(key, value []byte) error {
			fmt.Printf("key=%s, value=%s\n", key, value)
			return nil
		})
		return nil
	})
}

// // InsertUser inserts a user into the boltDB.
// func InsertUser(db *bolt.DB, username string, balance int, limit int, amount int, day int) error {
// 	u := User{balance, limit, amount, day}

// 	return db.Update(func(tx *bolt.Tx) error {
// 		bucket := tx.Bucket([]byte("users"))
// 		encodedUser, err := json.Marshal(u)
// 		if err != nil {
// 			return fmt.Errorf("error encoding User: %s", username)
// 		}
// 		return bucket.Put([]byte(username), encodedUser)
// 	})
// }
