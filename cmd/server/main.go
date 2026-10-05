package main

import (
	"fmt"

	"github.com/dbuck182/distribkv/internal/engine"
)

func main() {
	fmt.Println("distribkb initializing...")

	kv := engine.NewEngine()

	kv.Set("Name", []byte("Drew"))

	val := kv.Get("Nam")
	fmt.Printf("Name is: %s\n", string(val))

}
