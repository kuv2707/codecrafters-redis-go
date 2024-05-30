package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
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

func (s *Stream) appendEntry(id string, key string, value string) (string, error) {
	newid, err := validateId(id, s.entries)
	if err != nil {
		return "", err
	}
	s.entries = append(s.entries, StreamEntry{id: newid, data: key + string(0xff) + value})
	log(s.entries)
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
	log("processing ", parts, entries)
	if len(entries) == 0 {
		//this is the first id
		if parts[0] == "0" {
			if parts[1] == "*" {
				log("hereee")
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

func filterAsterisks(id string, store []StreamEntry) string {
	parts := strings.Split(id, "-")
	timestamp := strtoint(parts[0])
	replace := ""
	if len(store) > 0 {
		last := store[len(store)-1]
		lastparts := strings.Split(last.id, "-")
		lastts := strtoint(lastparts[0])
		if lastts == timestamp {
			replace = "0"
		}
	}
	if len(store) == 0 {
		replace = "1"
	} else {
		replace = "0"
	}
	parts[1] = strings.Replace(parts[1], "*", replace, 1)
	return strings.Join(parts, "-")
}
