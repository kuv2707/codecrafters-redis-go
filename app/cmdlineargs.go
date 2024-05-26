package main

import (
	"os"
	// "strings"
)

func parseCmdLineArgs() map[string]string {
	args := os.Args
	argmap := make(map[string]string)
	for i:=0;i<len(args);i++ {
		if args[i][:2] == "--" {
			argmap[args[i][2:]] = args[i+1]
		}
	}
	return argmap
}