package main

import (
	"fmt"
	"os"
)

var TEST = 07

func test() {
	if TEST == 0 {
		return
	}
	read := 0
	tc := []string{
		"*2",
		"$3",
		"GET",
		"$3",
		"foo",
	}
	res := ParseQuery(tc, &read)
	res.Print("")
	fmt.Println("->", read)

	res2 := ParseQuery(tc, &read)
	res2.Print("")
	fmt.Println("->", read)
	os.Exit(0)
}

func inc(a *int) {
	*a += 1
}
