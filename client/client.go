package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/apostolistselios/atm-simulator/api"
)

func main() {
	fmt.Println("WELCOME TO THE ATM.")
	userID, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}

	// Make the user request to the server.
	if err := verifyUser(userID); err != nil {
		log.Fatal(err)
	}

	// The main loop of the client starts.
	mainLoop(userID)
}

// getCredentials gets the userID from the user and returns it.
func getCredentials() (string, error) {
	var userID string
	fmt.Print("USER ID: ")
	if _, err := fmt.Scanf("%s\n", &userID); err != nil {
		return "", errors.New("error parsing the user ID")
	}
	return userID, nil
}

// verifyUser makes a HTTP request to /atm/{id} and handles the response.
func verifyUser(userID string) error {
	url := "http://localhost:8080/atm/" + userID
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	return nil
}

// mainLoop executes the main loop of the client.
func mainLoop(userID string) {
	answer := "y"
	for answer == "y" || answer == "Y" {
		userPrompt()
		var choice string
		if _, err := fmt.Scanf("%s\n", &choice); err != nil {
			log.Println("error incorrect choice")
			continue
		}

		choice = strings.ToLower(choice)
		if choice == "w" || choice == "d" {
			var amount int
			fmt.Print("PLEASE ENTER THE AMOUNT: ")
			if _, err := fmt.Scanf("%d\n", &amount); err != nil {
				log.Println("error incorrect amount")
				continue
			}

			// Check if the transaction is in the correct form.
			if err := checkTransaction(choice, amount); err != nil {
				log.Println(err)
				continue
			}

			// Make the transaction object.
			transaction := api.Transaction{
				UserID: userID,
				Type:   choice,
				Amount: amount,
			}

			// Make the transaction request to the server.
			if err := requestTransaction(transaction); err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("TRANSACTION COMPLETE")
		} else if choice == "b" {
			// Make the balance request to the server.
			balance, err := requestBalance(userID)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("YOUR BALANCE IS:", balance)
		} else if choice == "e" {
			fmt.Println("BYE BYE")
			break
		}
		// Check if the user wants to continue.
		fmt.Print("WOULD YOU LIKE TO CONTINUE (Y/N): ")
		fmt.Scanf("%s\n", &answer)
	}
}

func userPrompt() {
	fmt.Println("1. W TO WITHDRAW AN AMOUNT")
	fmt.Println("2. D TO DEPOSIT AN AMOUNT")
	fmt.Println("3. B TO SEE YOUR BALANCE")
	fmt.Println("4. E TO EXIT")
	fmt.Print("PLEASE CHOOSE THE TYPE A TYPE OF TRANSACTION: ")
}

// checkTransaction checks if the transaction is in the correct form.
// The transaction type (ttype) has to be w/W or d/D.
// The amount has to be a multiple of 20 or 50.
func checkTransaction(ttype string, amount int) error {
	if !(ttype == "w" || ttype == "W" || ttype == "d" || ttype == "D") {
		return errors.New("the transaction has to be between w/W or d/D, try again")
	}

	if amount%20 != 0 && amount%50 != 0 {
		return errors.New("the amount has to be multiple of 20 or multiple 50, try again")
	}
	return nil
}

// requestTransaction makes a HTTP request to /atm/{id}/transaction and handles the response.
func requestTransaction(transaction api.Transaction) error {
	// Encode the transaction object.
	encodedTrans, _ := json.Marshal(transaction)

	// Make the server request.
	url := "http://localhost:8080/atm/" + transaction.UserID + "/transaction"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(encodedTrans))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	return nil
}

// requestBalance makes a HTTP request to /atm/{id}/balance and handles the response.
func requestBalance(userID string) (string, error) {
	url := "http://localhost:8080/atm/" + userID + "/balance"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		return "", errors.New(string(data))
	}
	return string(data), nil
}
