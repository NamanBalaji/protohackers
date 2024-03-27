package main

import (
	"github.com/NamanBalaji/protohackers/4_unusual_database_program/db"
	"log"
	"net"
	"strings"
)

type Request struct {
	Key      string
	Value    string
	IsInsert bool
}

func GetRequest(data string) Request {
	key, value, found := strings.Cut(data, "=")

	return Request{
		Key:      key,
		Value:    value,
		IsInsert: found,
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Received error %s\n", err)
	}
	defer conn.Close()

	store := db.New("v1")

	for {
		buf := make([]byte, 1000)
		n, raddr, _ := conn.ReadFromUDP(buf)
		req := GetRequest(string(buf[0:n]))
		if req.IsInsert {
			store.Set(req.Key, req.Value)
		} else {
			res := store.Get(req.Key)
			conn.WriteToUDP([]byte(res), raddr)
		}
	}
}
