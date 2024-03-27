package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

const (
	upstreamAddr = "chat.protohackers.com:16963"
	tonyAddress  = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Received error %s\n", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Couldn't accept connection: %s", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	upstream, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		log.Fatalf("error dialing TCP address: %v", err)
	}
	defer upstream.Close()

	toConn, toUpstream := make(chan string, 10), make(chan string, 10)
	go Transform(bufio.NewReader(conn), toUpstream)
	go Transform(bufio.NewReader(upstream), toConn)

	for {
		select {
		case msg, ok := <-toConn:
			if !ok {
				return
			}
			io.WriteString(conn, msg)

		case msg, ok := <-toUpstream:
			if !ok {
				return
			}
			io.WriteString(upstream, msg)
		}
	}

}

var Regex = regexp.MustCompile(`^7[0-9a-zA-Z]{25,34}$`)

func Transform(in *bufio.Reader, ch chan string) {
	defer close(ch)

	for {
		msg, err := in.ReadString('\n')
		if err != nil {
			return
		}

		for _, word := range strings.Fields(msg) {
			if Regex.MatchString(word) {
				msg = strings.ReplaceAll(msg, word, tonyAddress)
			}
		}

		ch <- msg
	}
}
