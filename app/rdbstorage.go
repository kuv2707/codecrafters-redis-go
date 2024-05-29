package main

import (
	"fmt"
	"os"
	"time"
)

func testrdb() {
	// data := []byte{0b10110000, 0b01111111, 0b11111111, 0b11111111, 0b11111111} // 0b01000000 in binary, opcode will be 1
	// p := 0
	// result := parseLengthEncoded(data, &p)
	// log(result) // Outputs: 0
	// parseBytes()
	os.Exit(0)
}

func loadRDB(ctx *Context) {
	args := ctx.cmdArgs
	path := args["dir"] + "/" + args["dbfilename"]
	log(path)
	data, err := os.ReadFile(path)
	if err != nil {
		log(err)
		return
	}
	fmt.Printf("%x\n", data)
	ptr := 0
	parseBytes(data, &ptr, ctx)
	// return data
}

// 52 45 44 49 53 30 30 30 33 fa 0a 72 65 64 69 73 2d 62 69 74 73 c0 40 fa 09 72 65 64 69 73 2d 76 65 72 05 37 2e 32 2e 30 fe 00 fb 01 00 00 06 62 61 6e 61 6e 61 05 6d 61 6e 67 6f ff 2f ce 15 1d 67 ca 3c 2a 0a
// why 3 bytes after fb before string begins - it is the value type, 0 means value is string-encoded

// 52 45 44 49 53 30 30 30 33 fa 0a 72 65 64 69 73 2d 62 69 74 73 c0 40 fa 09 72 65 64 69 73 2d 76 65 72 05 37 2e 32 2e 30 fe 00 fb 04 00 00 09 70 69 6e 65 61 70 70 6c 65 09 62 6c 75 65 62 65 72 72 79 00 06 62 61 6e 61 6e 61 09 70 69 6e 65 61 70 70 6c 65 00 06 6f 72 61 6e 67 65 05 6d 61 6e 67 6f 00 09 62 6c 75 65 62 65 72 72 79 05 61 70 70 6c 65 ff 61 22 1d fd 84 0d 3a 64 0a

func parseBytes(data []byte, p *int, ctx *Context) {
	for data[*p] != 0xfb {
		*p += 1
	}
	*p += 1
	numKeys := parseLengthEncoded(data, p)
	*p += 1
	parseLengthEncoded(data, p)
	*p += 1
	totalKeys := numKeys
	for totalKeys > 0 {
		totalKeys -= 1
		key, val := parseKV(data, p)
		// *p += 1
		log("read - ", key, val)
		ctx.storage[key] = val
	}
	return

}

// 52 45 44 49 53 30 30 30 33 fa 09 72 65 64 69 73 2d 76 65 72 05 37 2e 32 2e 30 fa 0a 72 65 64 69 73 2d 62 69 74 73 c0 40 fe 00 fb 04 04 fc 00 9c ef 12 7e 01 00 00 00 09 72 61 73 70 62 65 72 72 79 06 6f 72 61 6e 67 65 fc 00 0c 28 8a c7 01 00 00 00 06 62 61 6e 61 6e 61 05 6d 61 6e 67 6f fc 00 0c 28 8a c7 01 00 00 00 05 67 72 61 70 65 09 72 61 73 70 62 65 72 72 79 fc 00 0c 28 8a c7 01 00 00 00 04 70 65 61 72 0a 73 74 72 61 77 62 65 72 72 79 ff 95 fc c9 2e 1b b4 b2 4e 0a

func parseKV(data []byte, p *int) (string, Value) {
	optionalopc := data[*p]
	log("opc", optionalopc)
	value := Value{expires: infiniteTime()}
	switch optionalopc {
	case 0xfc: // the expiry time is in little endian format
		{
			*p += 1
			value.expires = time.Time(time.UnixMilli(bytestoint64(data, p)))
		}
	case 0xfd:
		{
			*p += 1
			value.expires = time.Time(time.UnixMilli(1000 * int64(bytestoint(data, p))))
		}
	}
	*p += 1 // ignoring this byte as it stores info about type of value, which for now we know is string
	keylen := parseLengthEncoded(data, p)
	*p += 1
	key := string(data[*p : *p+keylen])
	*p += keylen
	vallen := parseLengthEncoded(data, p)
	*p += 1
	val := string(data[*p : *p+vallen])
	value.value = val
	*p += vallen
	return key, value
}

// pointer should point to the start of the length encoding
// increments the pointer to the last byte read
func parseLengthEncoded(data []byte, p *int) int {
	opcode := data[*p] >> 6
	switch opcode {
	case 0:
		return int((data[*p] << 2) >> 2)
	case 1:
		{

			sixbits := (int((data[*p]<<2)>>2) << 8)
			*p += 1
			next := data[*p]
			return int(sixbits) + int(next)
		}
	case 2:
		{
			num := bytestoint(data, p)
			return num
		}
	}
	return -1
}


// the following fns are for little endian
func bytestoint(data []byte, p *int) int {
	var num int = 0
	bitshift := 0
	for i := *p; i <= *p+3; i++ {
		num += int(data[i] << byte(bitshift))
		bitshift += 8
	}
	*p += 4
	return num
}
func bytestoint64(data []byte, p *int) int64 {
	var num int64 = 0
	bitshift := 0
	for i := *p; i <= *p+7; i++ {
		num += int64(data[i]) << bitshift
		bitshift += 8
	}
	*p += 8

	return num
}
