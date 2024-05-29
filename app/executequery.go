package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

var ackInfo = make(map[string]int)

func Execute(data *Data, conn net.Conn, ctx *Context, cmdctx *CommandContext) {
	// this data is a resp array, as per the protocol
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
						respond(conn, encodeBulkString(str))
					}
				case "PING":
					{
						respondIfMaster(ctx, conn, PONG)
					}
				case "SET":
					{
						key := data.children[i+1].content
						value := data.children[i+2].content
						dur := getDuration(data.children[i+3:])
						expires := time.Now().Add(dur)
						ctx.storage[key] = Value{
							value,
							expires,
						}
						respondIfMaster(ctx, conn, OK)
						propagateCommand(cmdctx, ctx)
					}
				case "GET":
					{
						key := data.children[i+1].content
						value, exists := ctx.storage[key]
						response := ""
						log(ctx.storage)
						if !exists || value.expired() {
							response = NULL_BULK_STRING
						} else {
							response = encodeBulkString(value.value)
						}
						respond(conn, response)
					}
				case "TYPE":
					{
						key := data.children[i+1].content
						value, exists := ctx.storage[key]
						response := ""
						log(ctx.storage)
						if !exists || value.expired() {
							response = encodeSimpleString("none")
						} else {
							response = encodeSimpleString("string")
						}
						respond(conn, response)
					}
				case "INFO":
					{
						respond(conn, encodeBulkString(replicationData(ctx)))
					}
				case "REPLCONF":
					{
						log("RECEIVED REPLCONF", ctx.offsetACK)
						subcomm := data.children[i+1].content
						switch strings.ToUpper(subcomm) {
						case "GETACK":
							respond(conn, encodeQuery("REPLCONF", "ACK", fmt.Sprint(ctx.offsetACK)))
						case "ACK":
							{
								slaveack := strtoint(data.children[i+2].content)
								laddr := conn.RemoteAddr().String()
								ackInfo[laddr] = slaveack
								ackUpdateChan <- AckUpdate{
									laddr:  laddr,
									ackVal: ackInfo[laddr],
								}
								log("Update ACK of", conn.RemoteAddr().String(), " to ", slaveack)
								// respond(conn, OK)
							}
						default:
							respond(conn, OK)
						}
					}
				case "PSYNC":
					{
						emptyRDB, _ := hex.DecodeString(EMPTY_RDB_HEX)
						byteslice := fmt.Sprintf("$%d\r\n%s", len(emptyRDB), string(emptyRDB))
						ctx.slaves[&conn] = true
						ackInfo[conn.RemoteAddr().String()] = 0
						log("Added slave", conn.RemoteAddr().String())
						respond(conn, encodeSimpleString(fmt.Sprintf("FULLRESYNC %s 0", ctx.info["master_replid"])))
						respond(conn, byteslice)
					}
				case "WAIT":
					{
						replNo := data.children[i+1].content
						timeout := data.children[i+2].content
						handleWait(strtoint(replNo), strtoint(timeout), conn, ctx)
						// respond(conn, encodeInteger(len(ctx.slaves)))
					}
				case "CONFIG":
					{
						subcomm := data.children[i+1].content
						switch strings.ToUpper(subcomm) {
						case "GET":
							{
								key := data.children[i+2].content
								respond(conn, encodeQuery(key, ctx.cmdArgs[key]))
							}
						}
					}
				case "KEYS":
					{
						keys := make([]string, 0, len(ctx.storage))
						for k := range ctx.storage {
							keys = append(keys, k)
						}
						respond(conn, encodeQuery(keys...))
					}
				}
				return // we need to return after processing this
			}
		case '*':
			{
				// maybe used when pipelining etc
				// recursively call Execute with this child
			}
		}
	}
	panic("Unhandled command")
}

type AckUpdate struct {
	laddr  string
	ackVal int
}

var ackUpdateChan = make(chan AckUpdate)

func handleWait(replNo int, timeout int, conn net.Conn, ctx *Context) {
	log("SERVER ACK is", ctx.offsetACK)
	if ctx.offsetACK == 0 {
		// obv slaves will also give 0, so dont even ask them
		respond(conn, encodeInteger(len(ctx.slaves)))
		return
	}

	// propagating this command increases server ACK by 37, but clients
	// will not acknowledge it till the next getack call.
	propagateCommand(&CommandContext{
		command: encodeQuery("REPLCONF", "GETACK", "*"),
		sender:  MASTER,
	}, ctx)
	valids := make(map[string]int)
	timeoutCh := time.After(time.Duration(timeout) * time.Millisecond)
	for true {
		select {
		case <-timeoutCh:
			{
				respond(conn, encodeInteger(len(valids)))
				return
			}
		case upd := <-ackUpdateChan:
			{

				fmt.Println("ACK update received", upd, ctx.offsetACK)
				if upd.ackVal+37 == ctx.offsetACK {
					valids[upd.laddr] = 1
				}
				if len(valids) == replNo {
					respond(conn, encodeInteger(replNo))
					return
				}
			}
		}

	}
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

func updateACKOffset(s string, ctx *Context) {
	log(ctx.info["role"], "Updating ACKOFF by ", len(s), " due to command ", s)
	ctx.offsetACK += len(s)
}

func propagateCommand(cmdctx *CommandContext, ctx *Context) {
	if ctx.info["role"] == "slave" {
		return
	}
	updateACKOffset(cmdctx.command, ctx)
	for slave := range ctx.slaves {
		log("sending to slave", (*slave).RemoteAddr())
		log("??", cmdctx.command, "??")
		(*slave).Write([]byte(cmdctx.command))
	}
}

func respondIfMaster(ctx *Context, conn net.Conn, res string) {
	if ctx.info["role"] == "master" {
		respond(conn, res)
	}
}

func respond(conn net.Conn, res string) {
	conn.Write([]byte(res))
}
