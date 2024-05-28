package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	// "time"
)

type Context struct {
	master    string
	info      map[string]string
	cmdArgs   map[string]string
	offsetACK int
	slaves    map[*net.Conn]bool
}

type SenderType = int

const CLIENT SenderType = 1
const MASTER SenderType = 2

type CommandContext struct {
	command string
	sender  SenderType
}

func main() {
	// test()
	args := parseCmdLineArgs()
	port := args["port"]
	if port == "" {
		port = "6379"
	}
	ctx := Context{master: "self", info: make(map[string]string), slaves: make(map[*net.Conn]bool), cmdArgs: args}
	ctx.info["port"] = port
	if args["replicaof"] != "" {
		ctx.info["role"] = "slave"
		ctx.info["master"] = args["replicaof"]
	} else {
		ctx.info["role"] = "master"
		ctx.info["master_replid"] = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		ctx.info["master_repl_offset"] = "0"
	}
	if ctx.info["role"] == "slave" {
		go connectToMaster(&ctx)

	}
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	fmt.Println("Listening on port " + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("Connection accepted successfully!")
		}
		go handleConnection(conn, &ctx, CLIENT)
		// l.Close()
	}
}

func connectToMaster(ctx *Context) {
	url := ctx.info["master"]
	url = func(url string) string {
		parts := strings.Split(url, " ")
		return parts[0] + ":" + parts[1]

	}(url)
	fmt.Printf("Connecting to master <%s>", url)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Failed to connect to master:", err)
		return
	}
	// defer conn.Close()
	conn.Write([]byte(encodeQuery("PING")))
	a := readConn(conn)
	log("master says", a)
	conn.Write([]byte(encodeQuery("REPLCONF", "listening-port", ctx.info["port"])))
	readConn(conn)
	log("master says", a)
	conn.Write([]byte(encodeQuery("REPLCONF", "capa", "psync2")))
	a = readConn(conn)
	conn.Write([]byte(encodeQuery("PSYNC", "?", "-1")))
	// readConn(conn)
	a = readConn(conn)
	log("master says", a)
	// parsePSYNCResponse(a)
	expectRDBFile(a, conn, ctx)
	handleConnection(conn, ctx, MASTER)

}

func expectRDBFile(a string, conn net.Conn, ctx *Context) {
	log("==", len(a), a)
	if len(a) > 56 {
		//rdb file is also present, so nothing to do
		a = a[56:]
	} else {
		a = readConn(conn)
	}
	// some commands may be present in the end of rdb file, now stored in a. Extract those commands and execute
	// log("cut a",len(a),a+"\n\n")
	spc := strings.Index(a, "\r\n")
	log(a[1:spc])
	size, _ := strconv.Atoi(a[1:spc])
	a = a[spc+2:]
	a = a[size:]
	if len(a) > 0 && a[0] == '*' {
		str := strings.Trim(a, "\r\n")
		commands := strings.Split(str, CRLF)
		processCommands(commands, conn, ctx, MASTER)
	}
	log("left part ->", a)
}

func readConn(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return ""
	}
	return string(buffer[:n])
}

func handleConnection(conn net.Conn, ctx *Context, sender SenderType) {
	// defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(ctx.info["role"], "Error reading:", err.Error())
			return
		}
		str := strings.Trim(string(buf[:n]), "\r\n")
		commands := strings.Split(str, CRLF)
		for _, c := range commands {
			log("->", c)
		}
		processCommands(commands, conn, ctx, sender)
	}
}


func processCommands(commands []string, conn net.Conn, ctx *Context, sender SenderType) {
	read := 0
	for true {
		read += processCommand(commands[read:], conn, ctx, sender)
		if read >= len(commands) {
			break
		}
	}
}
func processCommand(commandlist []string, conn net.Conn, ctx *Context, sender SenderType) int {

	r := 0
	data := ParseQuery(commandlist, &r)
	propag_cmd := strings.Join(commandlist[:r], "\r\n") + "\r\n" // alternatively, traverse the data and accumulate leaf node contents

	cmdctx := CommandContext{
		command: propag_cmd,
		sender:  sender,
	}
	// prop, respslv
	Execute(&data, conn, ctx, &cmdctx)
	if sender == MASTER {
		updateACKOffset(propag_cmd, ctx)
	}

	return r
}
