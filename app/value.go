package main

import (
	"time"
)

type Value struct {
	value   string
	expires time.Time
}

func (value *Value) expired() bool {
	return time.Now().After(value.expires)
}
