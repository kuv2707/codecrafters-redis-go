package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Context struct {
	role string
	master string
}

var TEST = 0

func main() {
	// test()
	args:=parseCmdLineArgs()
	port := args["port"]
	if port == "" {
		port = "6379"
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	}
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	ctx := Context{role:"master", master:"self"}
	if args["replicaof"] != "" {
		ctx.role = "slave"
		ctx.master = args["replicaof"]
	}
	fmt.Println(ctx)
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
		res := Execute(&data, ctx)
		conn.Write([]byte(res))
	}

}

func testQuery(comms []string) {
	// fmt.Println("testing:", comms)
	a, _ := ParseQuery(comms)
	// a.Print("")
	fmt.Println("->" + Execute(&a, nil))
}

func test() {
	if TEST != 0 {
		// // *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
		fmt.Println(time.Now().UnixNano())
		testQuery([]string{
			"*5",
			"$3",
			"SET",
			"$9",
			"blueberry",
			"$3",
			"BAR",
			"$2",
			"px",
			":2",
			"100",
		})
		fmt.Println(time.Now().UnixNano())
		time.Sleep(time.Duration(100) * time.Millisecond)
		testQuery([]string{
			"*2",
			"$3",
			"GET",
			"$3",
			"blueberry",
		})
		fmt.Println(time.Now().UnixNano())

		os.Exit(0)
	}
}