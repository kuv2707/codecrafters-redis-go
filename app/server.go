package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// "time"
)

type Context struct {
	master string
	info   map[string]string
}

var TEST = 0

func main() {
	args := parseCmdLineArgs()
	port := args["port"]
	if port == "" {
		port = "6379"
	}
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	ctx := Context{master: "self", info: make(map[string]string)}
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
		connectToMaster(&ctx)
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
		go handleConnection(conn, &ctx)
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
	defer conn.Close()
	conn.Write([]byte(encodeQuery("PING")))
	readConn(conn)
	conn.Write([]byte(encodeQuery("REPLCONF", "listening-port", ctx.info["port"])))
	readConn(conn)
	conn.Write([]byte(encodeQuery("REPLCONF", "capa", "psync2")))
	readConn(conn)
	conn.Write([]byte(encodeQuery("PSYNC", "?", "-1")))

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

func handleConnection(conn net.Conn, ctx *Context) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		str := string(buf[:len])
		data, _ := ParseQuery(strings.Split(str, "\r\n"))
		ress := Execute(&data, conn, ctx)
		for _, res := range ress {
			conn.Write([]byte(res))
		}
	}
}
