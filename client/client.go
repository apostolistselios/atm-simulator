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

	handleResponse(conn)
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
func handleResponse(conn net.Conn) {
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	fmt.Println(strings.Trim(response, "\n"))
}
