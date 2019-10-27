package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("WELCOME TO THE ATM.")
	username, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(conn, username)

	response := handleResponse(conn)
	if response != "OK" {
		log.Fatal(response)
	}

	answer := "y"
	for answer == "y" || answer == "Y" {
		fmt.Println("1. w/W <amount> to withdraw the amount")
		fmt.Println("2. d/D <amount> to deposit the amount")
		fmt.Print("Please choose the transaction you would like to do: ")

		var transaction string
		var amount int
		_, err := fmt.Scanf("%s %d\n", &transaction, &amount)
		if err != nil {
			log.Println("error", err)
			continue
		}

		err = checkTransaction(transaction, amount)
		if err != nil {
			log.Println(err)
			continue
		}

		// Send the message to the server. TODO: function sendmsg to server.
		msg := fmt.Sprintf("%s %d", transaction, amount)
		fmt.Fprintln(conn, msg)

		//TODO: Receive the server response.

		fmt.Print("Would you like to continue with another transaction (y/n): ")
		fmt.Scanf("%s\n", &answer)
	}
}

// getCredentials gets the username from the user and returns it.
func getCredentials() (string, error) {
	var username string

	fmt.Print("Username: ")
	_, err := fmt.Scanf("%s\n", &username)
	if err != nil {
		return "", fmt.Errorf("error parsing the username")
	}
	return username, nil
}

// handleResponse handles the response from the server.
func handleResponse(conn net.Conn) string {
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	return strings.Trim(response, "\n")
}

func checkTransaction(transaction string, amount int) error {
	if !(transaction == "w" || transaction == "W" || transaction == "d" || transaction == "D") {
		return fmt.Errorf("the transaction has to be between w/W or d/D, try again")
	}

	if amount%20 != 0 && amount%50 != 0 {
		return fmt.Errorf("the amount has to be multiple of 20 or multiple 50, try again")
	}
	return nil
}
