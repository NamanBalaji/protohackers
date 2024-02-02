package main

import (
	"bufio"
	"fmt"
	"github.com/NamanBalaji/protohackers/3_budget_chat/chat"
	"log"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Received error %s\n", err)
	}
	defer l.Close()

	chatroom := chat.New()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Couldn't accept connection: %s", err)
			continue
		}

		go handle(chatroom, conn)
	}
}

func handle(chatroom *chat.Chat, conn net.Conn) {
	fmt.Fprintln(conn, "Welcome to budgetchat! What shall I call you?")

	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(conn, err)
		conn.Close()
		return
	}
	username = strings.TrimSpace(username)

	if err := chatroom.AddUser(username, conn); err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	defer chatroom.RemoveUser(username)

	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		chatroom.SendMessage(username, strings.TrimSpace(in))
	}
}
