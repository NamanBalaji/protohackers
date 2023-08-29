package main

import (
	"io"
	"log"
	"net"
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
	a := conn.RemoteAddr().String()
	log.Printf("ACCEPT %s\n", a)

	written, err := io.Copy(conn, conn)
	if err != nil {
		log.Printf("ERROR %s %s\n", a, err)
	} else {
		log.Printf("CLOSE %s Wrote %d bytes\n", a, written)
	}
}
