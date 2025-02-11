package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type Stream struct {
	id      string
	entries []StreamEntry
}

type StreamEntry struct {
	id   string
	data []string
}

func getStream(key string, ctx *Context) *Stream {
	data, exists := ctx.storage[key]
	if !exists {
		return nil
	}
	if data.datatype == STREAM_TYPE {
		return data.value.(*Stream)
	}
	return nil
}

func createStream(id string) *Stream {
	return &Stream{
		id:      id,
		entries: make([]StreamEntry, 0, 10),
	}
}

func (s *Stream) appendEntry(id string, kvs ...string) (string, error) {
	newid, err := validateId(id, s.entries)
	if err != nil {
		return "", err
	}
	s.entries = append(s.entries, StreamEntry{id: newid, data: kvs})
	return newid, nil
}

func validateId(id string, entries []StreamEntry) (string, error) {
	if id == "0-0" {
		return "", errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	if id == "*" {
		return autoGenerateId(entries), nil
	}
	parts := strings.Split(id, "-")
	if len(entries) == 0 {
		//this is the first id
		if parts[0] == "0" {
			if parts[1] == "*" {
				return parts[0] + "-1", nil
			} else {
				return id, nil
			}
		} else {
			if parts[1] == "*" {
				return parts[0] + "-0", nil
			} else {
				return id, nil
			}
		}
	} else {
		lastparts := strings.Split(entries[len(entries)-1].id, "-")
		if isGreater(parts[0], lastparts[0]) {
			if parts[1] == "*" {
				// parts[0] is not 0, as it is greater than some no.
				return parts[0] + "-0", nil
			} else {
				return id, nil
			}
		} else if isEqual(parts[0], lastparts[0]) {
			if parts[1] == "*" {
				lastSeqNo := strtoint(lastparts[1])
				return parts[0] + "-" + fmt.Sprint(lastSeqNo+1), nil
			} else {
				if isGreater(parts[1], lastparts[1]) {
					return id, nil
				} else {
					return "", errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
				}
			}
		} else {
			return "", errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

}

func autoGenerateId(entries []StreamEntry) string {
	now := time.Now().UnixMilli()
	if len(entries) == 0 {
		return fmt.Sprintf("%d-0", now)
	} else {
		last := strings.Split(entries[len(entries)-1].id, "-")
		lastts := strtoint(last[0])
		if lastts == now {
			return fmt.Sprintf("%d-%d", now, strtoint(last[1])+1)
		} else {
			return fmt.Sprintf("%d-0", now)
		}
	}
}

func isGreater(first string, second string) bool {
	ms1 := strtoint(first)
	ms2 := strtoint(second)
	return ms1 > ms2
}
func isEqual(first string, second string) bool {
	ms1 := strtoint(first)
	ms2 := strtoint(second)
	return ms1 == ms2
}

func insertHyphen(id string, val int64) string {
	i := strings.Index(id, "-")
	if i == -1 {
		return id + "-" + fmt.Sprint(val)
	}
	return id
}

func getComps(id string) [2]int64 {
	a := strings.Split(id, "-")
	return [2]int64{strtoint(a[0]), strtoint(a[1])}
}

func xrange(s *Stream, start_entry_id string, end_entry_id string, ctx *Context) []StreamEntry {
	if start_entry_id == "-" {
		start_entry_id = s.entries[0].id
	}
	if end_entry_id == "+" {
		end_entry_id = s.entries[len(s.entries)-1].id
	}
	start_entry_id = insertHyphen(start_entry_id, 0)
	end_entry_id = insertHyphen(end_entry_id, math.MaxInt64)
	collected := make([]StreamEntry, 0)
	for i := range s.entries {
		if isIdInRangeInc(s.entries[i].id, start_entry_id, end_entry_id) {
			collected = append(collected, s.entries[i])
		}
	}
	return collected
}

func xread(args []string, ctx *Context) []Stream {
	l := len(args)
	collected := make([]Stream, 0)
	for i := 0; i < l/2; i++ {
		stream := getStream(args[i], ctx)
		selected := xrange(stream, justGreater(args[i+l/2]), "+", ctx)
		log(selected)
		if len(selected) > 0 {
			collected = append(collected, Stream{
				id:      stream.id,
				entries: selected,
			})
		}
	}
	return collected
}

func justGreater(id string) string {
	c := getComps(insertHyphen(id, 0))
	return fmt.Sprintf("%d-%d", c[0], c[1]+1)

}

func isIdInRangeInc(id string, low string, high string) bool {
	return isIdGT(id, low, true) && isIdGT(high, id, true)
}

func isIdGT(a string, b string, equality bool) bool {
	first := getComps(a)
	second := getComps(b)
	if first[0] > second[0] {
		return true
	} else if first[0] < second[0] {
		return false
	} else {
		if equality {
			return first[1] >= second[1]
		} else {
			return first[1] > second[1]
		}
	}
}

func processDollar(words []string, ctx *Context) []string {
	l := len(words)
	for i := l / 2; i < l; i++ {
		if words[i] == "$" {
			stream := getStream(words[i-l/2], ctx)
			if len(stream.entries) > 0 {
				words[i] = stream.entries[len(stream.entries)-1].id
			} else {
				log("NO ENTRIES YET")
			}
		}
	}
	log("AFTER $ PROCESS",words)
	return words
}
