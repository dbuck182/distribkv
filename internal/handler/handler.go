package handler

import (
	"log"

	"github.com/dbuck182/distribkv/internal/resp"
)

type CommandResponse struct {
	E       error
	Message string
}

func HandleCommand(commands resp.Value) CommandResponse {
	log.Printf("Read error: %b", commands.Typ)
	log.Printf("Hello")
	return CommandResponse{Message: "+OK\r\n"}
}
