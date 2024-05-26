package main

import "fmt"

const NULL_BULK_STRING = "$-1\r\n"

func encodeBulkString(s string) string {
	l := len(s)
	return fmt.Sprintf("$%d\r\n%s\r\n", l, s)
}

func encodeSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func encodeQuery(words []string) string{
	ret:=fmt.Sprintf("*%d\r\n",len(words))
	for _,word := range words {
		ret+=encodeBulkString(word)
	}
	return ret
}