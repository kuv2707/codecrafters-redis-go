package main

import (
	"math"
	"time"
)

type Value struct {
	value   string
	expires time.Time
}

func (value *Value) expired() bool {
	return time.Now().After(value.expires)
}

func nonExpireValue(val string) Value {
	return Value{
		value: val,
		expires: time.Now().Add(time.Duration(math.MaxInt64)),
	}
}