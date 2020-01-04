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
	ID      string `json:"userID"`
	Balance int    `json:"balance"`
	// Daily withdraw limit.
	WithdrawalLimit int `json:"withdrawalLimit"`
	// Amount withdrawed the current day.
	AmountWithdrawed int `json:"amountWithdrawed"`
	// The day the last withdraw transaction happened.
	LastWithdrawalDay int `json:"lastWithdrawalDay"`
}

// Transaction struct
type Transaction struct {
	UserID string `json:"userID"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
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
func ExecuteTransaction(db *bolt.DB, t Transaction) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(t.UserID))

		var u User
		json.Unmarshal(value, &u)
		if t.Type == "w" {
			if err := processWithdrawal(&u, &t); err != nil {
				return err
			}
		} else if t.Type == "d" {
			u.Balance += t.Amount
		}
		encodedUser, _ := json.Marshal(u)
		return bucket.Put([]byte(t.UserID), encodedUser)
	})
}

// processWithdrawal makes the required checks in order for the withdrawal to be carried on.
// If the checks pass the User is prepared and the withdrawal is ready to be executed.
// If an error occurs it returns the error, else it returns nil.
func processWithdrawal(u *User, t *Transaction) error {
	currDay := time.Now().Day()
	if u.Balance > 0 && u.Balance-t.Amount >= 0 {
		if currDay == u.LastWithdrawalDay && u.AmountWithdrawed+t.Amount > u.WithdrawalLimit {
			return fmt.Errorf("error withdrawal limit, Daily Limit: %d", u.WithdrawalLimit)
		} else if currDay == u.LastWithdrawalDay {
			u.Balance -= t.Amount
			u.AmountWithdrawed += t.Amount
		} else {
			u.Balance -= t.Amount
			u.AmountWithdrawed = t.Amount
			u.LastWithdrawalDay = currDay
		}
	} else {
		return fmt.Errorf("error high withdrawal amount, User: %s Balance: %d", t.UserID, u.Balance)
	}
	return nil
}

// CheckUserID checks if that specific userID exists in the database.
// If not returns an error.
func CheckUserID(db *bolt.DB, userID string) error {
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(userID))
		if value == nil {
			return fmt.Errorf("error User: %s doesn't exist", userID)
		}
		return nil
	})
	return err
}

// GetBalance returns the balance of the specific user.
func GetBalance(db *bolt.DB, userID string) (int, error) {
	var u User
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(userID))

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
// func InsertUser(db *bolt.DB, userID string, balance int, limit int, amount int, day int) error {
// 	u := User{userID, balance, limit, amount, day}

// 	return db.Update(func(tx *bolt.Tx) error {
// 		bucket, err := tx.CreateBucketIfNotExists([]byte("users"))
// 		encodedUser, err := json.Marshal(u)
// 		if err != nil {
// 			return fmt.Errorf("error encoding User: %s", userID)
// 		}
// 		return bucket.Put([]byte(u.ID), encodedUser)
// 	})
// }
