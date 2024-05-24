package main

import (
	"fmt"
	"strings"
)

var storage = make(map[string]string)
const NULL_BULK_STRING = "$-1\r\n"

func Execute(data *Data) string {
	// this data is an array, as per the protocol
	for i := 0; i < len(data.children); i++ {
		child := data.children[i].(Data)
		switch child.kind {
		case '$':
			{
				operator := child.children[0].(string)
				switch strings.ToUpper(operator) {
				case "ECHO":
					{
						str := data.children[i+1].(Data).children[0].(string)
						return encodeBlukString(str)
					}
				case "PING":
					{
						return "+PONG\r\n"
					}
				case "SET":
					{
						key := data.children[i+1].(Data).children[0].(string)
						value := data.children[i+2].(Data).children[0].(string)
						storage[key] = value
						return encodeSimpleString("OK")
					}
				case "GET":
					{
						key := data.children[i+1].(Data).children[0].(string)
						value, exists := storage[key]
						if !exists {
							return NULL_BULK_STRING
						}
						return encodeBlukString(value)
					}
				}
			}
		}
	}
	return "null"
}

func encodeBlukString(s string) string {
	l := len(s)
	return fmt.Sprintf("$%d\r\n%s\r\n", l, s)
}

func encodeSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n",s)
}