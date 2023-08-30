package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net"
)

type request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("recieved error %s\n", err)
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

	buf := bufio.NewReader(conn)
	for {
		req, err := buf.ReadBytes('\n')
		if err != nil {
			log.Printf("an error occured while reading request %v\n", err)

			return
		}

		resp, err := handleRequest(req)
		if err != nil {
			log.Println(err)
			conn.Write([]byte(err.Error()))

			return
		}

		conn.Write(append(resp, []byte("\n")...))
	}
}

func handleRequest(inp []byte) ([]byte, error) {
	var req request
	if err := json.Unmarshal(inp, &req); err != nil {
		return nil, err
	}

	if req.Method == nil || req.Number == nil {
		return nil, errors.New("bad request")
	}

	if *req.Method != "isPrime" {
		return nil, errors.New("method not supported")
	}

	return json.Marshal(response{
		Method: *req.Method,
		Prime:  isPrime(*req.Number),
	})
}

func isPrime(value float64) bool {
	if float64(int64(value)) != value {
		return false
	}

	for i := 2.0; i <= math.Sqrt(value); i++ {
		if math.Mod(value, i) == 0 {
			return false
		}
	}

	return value > 1
}
