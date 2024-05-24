package main

import "fmt"

type Data struct {
	kind     byte
	size     int
	children []interface{}
}

var TAB string = "   "

func (data *Data) Print(spc string) {
	// return
	// fmt.Println("->",data.kind)
	switch data.kind {
	case '*':
		fmt.Println(spc+"[")
		for _,child := range data.children {
			dat := child.(Data)
			dat.Print(spc + TAB)
		}
		fmt.Println(spc+"]")
	case '$':
		str := data.children[0].(string)
		fmt.Println(spc+str)
	}
}
