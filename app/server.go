package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var TEST = 0

func main() {
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
	// fmt.Println("testing:", comms)
	a, _ := ParseQuery(comms)
	// a.Print("")
	fmt.Println("->" + Execute(&a))
}
