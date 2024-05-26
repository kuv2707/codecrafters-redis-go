package main

import (
	// "fmt"
	"strconv"
	// "strings"
)

func parseHeader(header string) (byte, int) {
	kind := header[0]
	size, _ := strconv.Atoi(header[1:])
	return kind, size

}

func ParseQuery(query []string, read *int) Data {
	log(*read)
	kind, size := parseHeader(query[*read])
	switch kind {
	case '*':
		return parseArray(query, size, read)
	case '$':
		return parseBulkString(query, size, read)
	case ':':
		return parseInteger(query, size, read)
	}

	return Data{}
}

func parseArray(query []string, size int, read *int) Data {
	*read += 1
	data := Data{kind: '*', size: size, children: make([]Data, size)}
	for i := 0; i < size; i++ {
		elem := ParseQuery(query, read)
		data.children[i] = elem
	}
	return data
}

func parseBulkString(query []string, size int, read *int) Data {
	data := Data{kind: '$', size: size, children: make([]Data, 1)}
	data.content = query[*read+1]
	*read += 2
	return data
}

func parseInteger(query []string, size int, read *int) Data {
	data := Data{kind: ':', size: size, children: make([]Data, 1)}
	data.content = query[1]
	*read += 2
	return data
}
