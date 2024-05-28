package main

import (
	"fmt"
	"strconv"
)

func log(a ...interface{}) {
	fmt.Println(a...)
}

func strtoint(s string) int {
	size, _ := strconv.Atoi(s)
	return size
}