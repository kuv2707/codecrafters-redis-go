package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

var storage = make(map[string]Value)

var slaves = make(map[*net.Conn]bool)

func Execute(data *Data, conn net.Conn, ctx *Context) ([]string, bool) {
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
						return []string{encodeBulkString(str)}, false
					}
				case "PING":
					{
						return []string{PONG},false
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
						return []string{encodeSimpleString("OK")}, true
					}
				case "GET":
					{
						key := data.children[i+1].content
						value, exists := storage[key]
						if !exists {
							return []string{NULL_BULK_STRING}, false
						}
						if value.expired() {
							return []string{NULL_BULK_STRING}, false
						}
						return []string{encodeBulkString(value.value)},false
					}
				case "INFO":
					return []string{encodeBulkString(replicationData(ctx))}, false
				case "REPLCONF":
					return []string{OK}, false
				case "PSYNC":
					emptyRDB, _ := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
					byteslice := fmt.Sprintf("$%d\r\n%s", len(emptyRDB), string(emptyRDB))
					slaves[&conn] = true
					return []string{encodeSimpleString(fmt.Sprintf("FULLRESYNC %s 0", ctx.info["master_replid"])), byteslice}, false
				}
			}
		}
	}
	return []string{"null"}, false
}

func replicationData(ctx *Context) string {
	ret := ""
	for k, v := range ctx.info {
		ret += k + ":" + v + "\n"
	}
	return ret
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
