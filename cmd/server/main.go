package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/dbuck182/distribkv/internal/engine"
	"github.com/dbuck182/distribkv/internal/resp"
)

func main() {
	portPtr := flag.String("port", "8080", ":<port#>")
	flag.Parse()
	fmt.Println("distribkb initializing... on port", *portPtr)
	kv := engine.NewEngine()
	// Okay I need to create a loop that will
	// start the tcp server and accept connections
	port := ":" + *portPtr
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server retry: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection Error at: %v", time.Now())
		}

		go handleConnection(conn)
	}

	kv.Set("Name", []byte("Drew"))

	val := kv.Get("Nam")
	fmt.Printf("Name is: %s\n", string(val))

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := resp.NewReader(conn)
	for {
		message, err := reader.Read()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		msg := strings.ToUpper(strings.TrimSpace(message))
		response := fmt.Sprintf("ACK: %s\n", msg)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}
