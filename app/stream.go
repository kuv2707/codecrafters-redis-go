package main

import (
	"errors"
	"strings"
)

type Stream struct {
	id      string
	entries []StreamEntry
}

type StreamEntry struct {
	id   string
	data string
}

func createStream(id string) *Stream {
	return &Stream{
		id:      id,
		entries: make([]StreamEntry, 0, 10),
	}
}

func (s *Stream) appendEntry(id string, key string, value string) error {
	if id == "0-0" {
		return errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	if len(s.entries) > 0 {
		lastEntry := s.entries[len(s.entries)-1]
		if compareID(id, lastEntry.id) != 1 {
			return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}
	s.entries = append(s.entries, StreamEntry{id: id, data: key + string(0xff) + value})
	log(s.entries)
	return nil
}

func compareID(first string, second string) int {
	firstparts := strings.Split(first, "-")
	secondparts := strings.Split(second, "-")
	log("comparing", firstparts, secondparts)
	ms1 := strtoint(firstparts[0])
	ms2 := strtoint(secondparts[0])
	if ms1 > ms2 {
		return 1
	} else if ms1 < ms2 {
		return -1
	} else {
		se1 := strtoint(firstparts[1])
		se2 := strtoint(secondparts[1])
		if se1 > se2 {
			return 1
		} else if se1 < se2 {
			return -1
		} else {
			return 0
		}
	}
}
