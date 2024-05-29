package main

import "fmt"

const NULL_BULK_STRING = "$-1\r\n"
const PONG = "+PONG\r\n"
const OK = "+OK\r\n"
const CRLF = "\r\n"
const EMPTY_RDB_HEX = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

func encodeBulkString(s string) string {
	l := len(s)
	return fmt.Sprintf("$%d\r\n%s\r\n", l, s)
}

func encodeSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func encodeInteger(a int) string {
	return fmt.Sprintf(":%d\r\n", a)
}

// encodes the words passed and composes a query (array with encoded words)
func encodeQuery(words ...string) string {
	ret := fmt.Sprintf("*%d\r\n", len(words))
	for _, word := range words {
		ret += encodeBulkString(word)
	}
	return ret
}
