package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// User struct
type User struct {
	Balance int
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

// InsertUser inserts a user into the boltDB.
func InsertUser(db *bolt.DB, username string, balance int) error {
	u := User{balance}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		encodedUser, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("error encoding User: %s", username)
		}
		return bucket.Put([]byte(username), encodedUser)
	})
}

// UpdateBalance updates the balance of the user with that specific username.
// transaction slice [username, tranType, amount]
func UpdateBalance(db *bolt.DB, transaction Transaction) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(transaction.UserID))
		if value == nil {
			return fmt.Errorf("error User: %s doesn't exist", transaction.UserID)
		}

		var u User
		err := json.Unmarshal(value, &u)
		if err != nil {
			return fmt.Errorf("error decoding User: %s", transaction.UserID)
		}

		if transaction.Type == "w" {
			if u.Balance > 0 && u.Balance-transaction.Amount >= 0 {
				u.Balance -= transaction.Amount
			} else {
				return fmt.Errorf("error withdrawal amount too high, User: %s Balance: %d", transaction.UserID, u.Balance)
			}
		} else if transaction.Type == "d" {
			u.Balance += transaction.Amount
		}

		encodedUser, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("error encoding User: %s", transaction.UserID)
		}
		return bucket.Put([]byte(transaction.UserID), encodedUser)
	})
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

		err := json.Unmarshal(value, &u)
		if err != nil {
			return fmt.Errorf("error decoding User: %s", username)
		}
		return nil
	})

	return u, err
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
