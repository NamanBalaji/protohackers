package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listner, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Printf("Error listening: %s\n", err)
		os.Exit(1)
	}

	defer listner.Close()

	fmt.Printf("Server listening on port: %d\n", 9000)

	for {
		client, err := listner.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %s\n", err)
			continue
		}

		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	fmt.Printf("Connected: %s\n", client.RemoteAddr())
	defer client.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := client.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading: %s\n", err)
			return
		}

		_, err = client.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
			return
		}
	}
}
