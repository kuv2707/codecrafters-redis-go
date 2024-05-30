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


func mapDataArrayToContent(data []Data) []string {
	strs := make([]string, 0, 10)
	for a := range data {
		strs = append(strs, data[a].content)
	}
	return strs
}