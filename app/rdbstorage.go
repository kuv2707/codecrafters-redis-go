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

func loadRDB(ctx *Context) string {
	args := ctx.cmdArgs
	path := args["dir"] + "/" + args["dbfilename"]
	log(path)
	data, err := os.ReadFile(path)
	if err != nil {
		log(err)
		return ""
	}
	fmt.Printf("%x\n", data)
	ptr := 0
	return parseBytes(data, &ptr)
	// return data
}

// 52 45 44 49 53 30 30 30 33 fa 0a 72 65 64 69 73 2d 62 69 74 73 c0 40 fa 09 72 65 64 69 73 2d 76 65 72 05 37 2e 32 2e 30 fe 00 fb 01 00 00 06 62 61 6e 61 6e 61 05 6d 61 6e 67 6f ff 2f ce 15 1d 67 ca 3c 2a 0a 

func parseBytes(data []byte, p *int) string {
	for data[*p] != 0xfb {
		*p += 1
	}
	*p += 1
	parseLengthEncoded(data, p)
	*p += 1
	parseLengthEncoded(data, p)
	*p += 1
	parseLengthEncoded(data, p)
	*p += 1
	strlen := parseLengthEncoded(data, p)
	*p += 1
	str := string(data[*p : *p+strlen])
	return str

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
