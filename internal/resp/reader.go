package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// This is the deserializer and should read raw
// bytes from a connection and parse them into Go values

type Reader struct {
	rd *bufio.Reader
}

type Value struct {
	Typ   byte
	Str   string
	Num   int64
	Bulk  string
	Array []Value
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(rd),
	}
}

func (r *Reader) Read() (response Value, err error) {

	message, err := r.rd.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch message {
	case '+':
		return r.readSimpleString()
	case '$':
		return r.readBulkString()
	case '*':
		return r.readArray()
	case '-':
		return r.readError()
	case ':':
		return r.readInteger()
	default:
		return Value{}, fmt.Errorf("Unknown resp prefix: %q", message)
	}

}

func (r *Reader) readArray() (response Value, err error) {
	message, err := r.rd.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}
	loops, err := strconv.Atoi(string(message))
	arr_vals := Value{Typ: '*', Array: make([]Value, 0, loops)}
	for index := range loops {
		// Go through each section
		message, err := r.rd.ReadByte()
		if err != nil {
			return Value{}, err
		}
		switch message {
		case '+':
			resp, err := r.readSimpleString()
			if err != nil {
				arr_vals.Array[index] = resp
			}
		case '$':
			resp, err := r.readBulkString()
			if err != nil {
				arr_vals.Array[index] = resp
			}
		case '*':
			resp, err := r.readArray()
			if err != nil {
				arr_vals.Array[index] = resp
			}
		case '-':
			resp, err := r.readError()
			if err != nil {
				arr_vals.Array[index] = resp
			}
		case ':':
			resp, err := r.readInteger()
			if err != nil {
				arr_vals.Array[index] = resp
			}
		default:
			return Value{}, fmt.Errorf("Unknown resp prefix: %q", message)
		}

	}
	return arr_vals, nil
}

func (r *Reader) readBulkString() (response Value, err error) {
	// First have to readline in order to get the length of the following
	// Then we can use the read certain number of bytes to get our string
	// set it equal to the bulkstring and return

	bytes, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	bytes_num, err := strconv.Atoi(bytes)
	if err != nil {
		return Value{}, err
	}
	// Now we will read until we get that number of bytes
	buf := make([]byte, bytes_num)

	_, err = io.ReadFull(r.rd, buf)
	if err != nil {
		return Value{}, err
	}
	// Now strip the crlf and all set
	str, err := stripCRLF(string(buf))
	if err != nil {
		return Value{}, err
	}
	return Value{Typ: '$', Bulk: str}, nil
}

func (r *Reader) readError() (response Value, err error) {
	trimmed_message, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	resp := Value{Typ: '-', Str: trimmed_message}
	return resp, nil
}

func (r *Reader) readInteger() (response Value, err error) {
	trimmed_message, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	resp := Value{Typ: ':', Str: trimmed_message}
	return resp, nil
}

func (r *Reader) readSimpleString() (response Value, err error) {

	trimmed_message, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	resp := Value{Typ: '+', Str: trimmed_message}
	return resp, nil
}

func (r *Reader) readLine() (string, error) {
	message, err := r.rd.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmed_message, err := stripCRLF(message)
	if err != nil {
		return "", err
	}
	return trimmed_message, nil
}

func stripCRLF(passed string) (returned string, err error) {
	if len(passed) < 2 || passed[len(passed)-2] != '\r' {
		return "", fmt.Errorf("Invalid RESP CRLF termination")
	}
	return passed[:len(passed)-2], nil
}
