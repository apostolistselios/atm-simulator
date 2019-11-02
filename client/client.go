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
	username, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8080/atm/user"
	resp, err := http.Post(url, "text/plain", strings.NewReader(username))
	if err != nil {
		log.Println(err)
	}
	err = handleResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	answer := "y"
	for answer == "y" || answer == "Y" {
		fmt.Println("1. W <AMOUNT> TO WITHDRAW THE AMOUNT")
		fmt.Println("2. D <AMOUNT> TO DEPOSIT THE AMOUNT")
		fmt.Print("PLEASE CHOOSE THE TYPE A TYPE OF TRANSACTION: ")

		var tranType string
		var amount int
		_, err := fmt.Scanf("%s %d\n", &tranType, &amount)
		if err != nil {
			log.Println("error incorrect transaction form")
			continue
		}

		// Check if the transaction is in the correct form.
		err = checkTransaction(tranType, amount)
		if err != nil {
			log.Println(err)
			continue
		}

		// Make the transaction object.
		transaction := api.Transaction{
			UserID: username,
			Type:   tranType,
			Amount: amount,
		}

		// Encode the transaction object.
		buffer := new(bytes.Buffer)
		err = json.NewEncoder(buffer).Encode(transaction)
		if err != nil {
			log.Println(err)
			continue
		}

		// Make the server request.
		url := "http://localhost:8080/atm/user/transaction"
		resp, err := http.Post(url, "application/json", buffer)
		if err != nil {
			log.Println(err)
		}

		// Handle the server response.
		err = handleResponse(resp)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("TRANSACTION COMPLETE")
		}

		// Check if the user wants to continue.
		fmt.Print("WOULD YOU LIKE TO CONTINUE OR EXIT (Y/N): ")
		fmt.Scanf("%s\n", &answer)
	}
}

// getCredentials gets the username from the user and returns it.
func getCredentials() (string, error) {
	var username string

	fmt.Print("Username: ")
	_, err := fmt.Scanf("%s\n", &username)
	if err != nil {
		return "", errors.New("error parsing the username")
	}
	return username, nil
}

// handleResponse handles a server response.
func handleResponse(resp *http.Response) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("error reading response body")
	}

	if string(body) != "OK" {
		return errors.New(string(body))
	}
	return nil
}

// checkTransaction checks if the transaction is in the correct form.
// The transaction type (tranType) has to be w/W or d/D.
// The amount has to be a multiple of 20 or 50.
func checkTransaction(tranType string, amount int) error {
	if !(tranType == "w" || tranType == "W" || tranType == "d" || tranType == "D") {
		return errors.New("the transaction has to be between w/W or d/D, try again")
	}

	if amount%20 != 0 && amount%50 != 0 {
		return errors.New("the amount has to be multiple of 20 or multiple 50, try again")
	}
	return nil
}
