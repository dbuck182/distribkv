package resp

import (
	"io"
)

// This is the serializer which takes Go values and
// formats them back into raw RESP bytes for the client

// This function transforms the message
func WriteToResp(writer io.Writer, output Value) {
	var finalMessage string
	switch output.Typ {
	case '+':
		finalMessage = writeSimpleString(output.Str)
	case '$':
		finalMessage = writeBulkString(output.Bulk)
	case '*':
		finalMessage = writeArray()
	case '-':
		finalMessage = writeError(output.Str)
	case ':':
		finalMessage = writeInteger(output.Num)
	default:
		finalMessage = "-Error with output message\r\n"
	}

	writer.Write([]byte(finalMessage))
}

func writeSimpleString(message string) string {
	return "+" + message + "\r\n"
}

func writeError(message string) string {
	return "-" + message + "\r\n"
}

func writeInteger(num int64) string {
	return ":" + string(num) + "\r\n"
}

func writeBulkString(message string) string {
	// Add the length to the front
	length := len(message)
	return string(length) + "\r\n" + message + "\r\n"
}

func writeArray() string {
	return "-ArraySorry\r\n"
}
