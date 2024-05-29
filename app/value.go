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
	log("expires at ",value.expires)
	return time.Now().After(value.expires)
}

func nonExpireValue(val string) Value {
	return Value{
		value: val,
		expires: time.Now().Add(time.Duration(math.MaxInt64)),
	}
}

func expireValue(val string, exp time.Time) Value {
	return Value{
		value: val,
		expires: exp,
	}
}

func infiniteTime() time.Time {
	return time.Now().Add(time.Duration(math.MaxInt64))
}