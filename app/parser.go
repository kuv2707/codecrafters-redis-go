package main

import (
	"fmt"
	"strconv"
	// "strings"
)

func parseHeader(header string) (byte, int) {
	kind := header[0]
	size, _ := strconv.Atoi(header[1:])
	return kind, size

}

func ParseQuery(query []string) (Data, int) {
	kind, size := parseHeader(query[0])
	fmt.Println(query)
	switch kind {
	case '*':
		return parseArray(query[1:], size)
	case '$':
		return parseBulkString(query, size)
	}

	return Data{}, 0
}

func parseArray(query []string, size int) (Data, int) {
	ind := 0
	data := Data{kind: '*', size: size, children: make([]interface{}, size)}
	for i := 0; i < size; i++ {
		elem, off := ParseQuery(query[ind:])
		// fmt.Println(elem)
		ind += off
		data.children[i] = elem
	}
	return data, ind
}

func parseBulkString(query []string, size int) (Data, int) {
	data := Data{kind: '$', size: size, children: make([]interface{}, 1)}
	data.children[0] = query[1]
	return data, 2
}
