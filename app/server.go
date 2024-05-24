package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var TEST = 0

func main() {
	if TEST != 0 {
		// // *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
		testQuery([]string{
			"*3",
			"$3",
			"SET",
			"$3",
			"FOO",
			"$3",
			"BAR",
		})
		testQuery([]string{
			"*2",
			"$3",
			"GET",
			"$3",
			"FOO",
		})
		return
	}

	l, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening on port 6379")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("Connection accepted successfully!")
		}
		go handleConnection(conn)
		// l.Close()
	}
}

func handleConnection(conn net.Conn) {
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
		res := Execute(&data)
		conn.Write([]byte(res))
	}

}

func testQuery(comms []string) {
	fmt.Println("testing:", comms)
	a, _ := ParseQuery(comms)
	a.Print("")
	fmt.Println("->" + Execute(&a))
}
