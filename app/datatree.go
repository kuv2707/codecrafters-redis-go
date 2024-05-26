package main

import (
	"fmt"
	"strconv"
	// "strconv"
)

type Data struct {
	kind     byte
	size     int
	content  string
	children []Data
}

var TAB string = "   "

func (data *Data) Print(spc string) {
	// return
	// fmt.Println("->",data.kind)
	switch data.kind {
	case '*':
		fmt.Println(spc+"[")
		for _,child := range data.children {
			dat := child
			dat.Print(spc + TAB)
		}
		fmt.Println(spc+"]")
	case '$':
		str := data.content
		fmt.Println(spc+str)
	case ':':
		str := data.content
		fmt.Println(spc+str)
	}
}

func (data *Data) asInt() (int, bool) {
    // if data.kind != ':' {
    //     return 0, false
    // }
	fmt.Println("converting ",data.content)
    num, err := strconv.Atoi(data.content)
    if err != nil {
        return 0, false
    }
    return num, true
}