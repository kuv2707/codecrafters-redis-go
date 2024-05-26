package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var storage = make(map[string]Value)

func Execute(data *Data, ctx *Context) string {
	// this data is an array, as per the protocol
	for i := 0; i < len(data.children); i++ {
		child := data.children[i]
		switch child.kind {
		case '$':
			{
				operator := child.content
				switch strings.ToUpper(operator) {
				case "ECHO":
					{
						str := data.children[i+1].content
						return encodeBulkString(str)
					}
				case "PING":
					{
						return "+PONG\r\n"
					}
				case "SET":
					{
						key := data.children[i+1].content
						value := data.children[i+2].content
						dur := getDuration(data.children[i+3:])
						expires := time.Now().Add(dur)
						fmt.Println("set at ", time.Now().UnixMicro())
						storage[key] = Value{
							value,
							expires,
						}
						return encodeSimpleString("OK")
					}
				case "GET":
					{
						key := data.children[i+1].content
						value, exists := storage[key]
						if !exists {
							return NULL_BULK_STRING
						}
						if value.expired() {
							return NULL_BULK_STRING
						}
						return encodeBulkString(value.value)
					}
				case "INFO":
					return encodeBulkString("role:"+ctx.role)
				}
			}
		}
	}
	return "null"
}

func encodeBulkString(s string) string {
	l := len(s)
	return fmt.Sprintf("$%d\r\n%s\r\n", l, s)
}

func encodeSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func getDuration(data []Data) time.Duration {
	if len(data) < 2 {
		return time.Duration(math.MaxInt64)
	}
	if data[0].kind == '$' && strings.EqualFold(data[0].content, "PX") {
		num, valid := data[1].asInt()
		if valid {
			return time.Duration(num) * time.Millisecond
		}
	}
	return time.Duration(math.MaxInt64)
}
