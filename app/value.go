package main

import (
	"time"
)

type Value struct {
	value   string
	expires time.Time
}

const NULL_BULK_STRING = "$-1\r\n"

func (value *Value) expired() bool {
	return time.Now().After(value.expires)
}
