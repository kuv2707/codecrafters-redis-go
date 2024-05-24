package main

import (
	"fmt"
	"strings"
)

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
						str := data.children[1].(Data).children[0].(string)
						return encodeBlukString(str)
					}
				case "PING":
					{
						return "+PONG\r\n"
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
