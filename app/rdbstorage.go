package main

import (
	"fmt"
	"os"
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
// why 3 bytes after fb before string begins


// 52 45 44 49 53 30 30 30 33 fa 0a 72 65 64 69 73 2d 62 69 74 73 c0 40 fa 09 72 65 64 69 73 2d 76 65 72 05 37 2e 32 2e 30 fe 00 fb 04 00 00 09 70 69 6e 65 61 70 70 6c 65 09 62 6c 75 65 62 65 72 72 79 00 06 62 61 6e 61 6e 61 09 70 69 6e 65 61 70 70 6c 65 00 06 6f 72 61 6e 67 65 05 6d 61 6e 67 6f 00 09 62 6c 75 65 62 65 72 72 79 05 61 70 70 6c 65 ff 61 22 1d fd 84 0d 3a 64 0a

func parseBytes(data []byte, p *int, ctx *Context) {
	for data[*p] != 0xfb {
		*p += 1
	}
	*p += 1
	numKeys := parseLengthEncoded(data, p)
	*p += 1
	parseLengthEncoded(data, p)
	*p += 2
	for numKeys > 0 {
		numKeys -= 1
		key, val := parseKV(data, p)
		*p+=1
		log("read - ",key,val)
		ctx.storage[key] = nonExpireValue(val)

	}
	return

}

func parseKV(data []byte, p *int) (string, string) {
	keylen := parseLengthEncoded(data, p)
	*p += 1
	key := string(data[*p : *p+keylen])
	*p += keylen
	vallen := parseLengthEncoded(data, p)
	*p += 1
	val := string(data[*p : *p+vallen])
	*p += vallen
	return key, val
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
			num := int(data[*p+1])<<24 + int(data[*p+2])<<16 + int(data[*p+3])<<8 + int(data[*p+4])
			*p += 4
			return num
		}
	}
	return -1
}
