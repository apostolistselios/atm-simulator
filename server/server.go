package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on :8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(message)

		conn.Write([]byte("hello back\n"))
	}
}
