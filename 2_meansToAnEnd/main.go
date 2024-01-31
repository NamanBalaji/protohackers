package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Message struct {
	Type uint8
	Arg1 int32
	Arg2 int32
}

type Record struct {
	Timestamp int32
	Price     int32
}

type Store struct {
	records []Record
}

func decodePayload(data []byte) Message {
	return Message{
		Type: data[0],
		Arg1: int32(binary.BigEndian.Uint32(data[1:5])),
		Arg2: int32(binary.BigEndian.Uint32(data[5:])),
	}
}

func (t *Store) insert(m Message) {
	var r Record
	r.Price = m.Arg2
	r.Timestamp = m.Arg1
	t.records = append(t.records, r)
}

func (t *Store) query(m Message) int32 {
	minTime := m.Arg1
	maxTime := m.Arg2

	var n int
	var total int
	for _, r := range t.records {
		if r.Timestamp >= minTime && r.Timestamp <= maxTime {
			n++
			total += int(r.Price)
		}
	}

	if n == 0 {
		return 0
	}

	return int32(float64(total) / float64(n))
}

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
	clientStore := &Store{}

	buf := make([]byte, 9)

	for {
		_, err := io.ReadAtLeast(conn, buf, 9)
		if err != nil {
			log.Printf("conn.Read(); %v", err)
			return
		}

		msg := decodePayload(buf)

		switch msg.Type {
		case 'I':
			clientStore.insert(msg)
		case 'Q':
			mean := clientStore.query(msg)
			if err := binary.Write(conn, binary.BigEndian, mean); err != nil {
				log.Println("error in writing response ", err)
			}
		}
	}
}
