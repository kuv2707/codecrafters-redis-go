package main

import (
	"fmt"
	"strconv"
)

func log(a ...interface{}) {
	fmt.Println(a...)
}

func strtoint(s string) int64 {
	size, _ := strconv.ParseInt(s,10, 64)
	return size
}