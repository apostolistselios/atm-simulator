package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Welcome to the ATM.")
	username, password, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(username, password)

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	fmt.Println(response)
}

// getCredentials gets from the user the username and password and returns them.
func getCredentials() (string, string, error) {
	var username string
	var password string
	fmt.Print("Username: ")
	_, err := fmt.Scanf("%s\n", &username)
	if err != nil {
		return "", "", err
	}
	fmt.Print("Password: ")
	_, err = fmt.Scanf("%s\n", &password)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
