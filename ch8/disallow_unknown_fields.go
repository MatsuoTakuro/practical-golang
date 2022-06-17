package ch8

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Rectangle struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func disallowUnkownFields() {
	f, err := os.Open("./ch8/radius.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var rect Rectangle
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	if err := d.Decode(&rect); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", rect)
}
