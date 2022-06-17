package ch8

import (
	"encoding/json"
	"fmt"
)

func encodeSlice() {
	langsNull()
	langsEmptyArray()

}

func langsNull() {
	u := user{
		UserID:   "001",
		UserName: "gopher",
	}
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
}

func langsEmptyArray() {
	u := user{
		UserID:    "001",
		UserName:  "gopher",
		Languages: []string{},
	}
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
}
