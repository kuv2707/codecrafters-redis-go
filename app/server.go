package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	defer l.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}else {
		fmt.Println("Connection accepted successfully!")
	}
	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf:=make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	str := string(buf[:len])
	if strings.Contains(str, "PING") {
		conn.Write([]byte("+PONG\r\n"))
	}
	
}
