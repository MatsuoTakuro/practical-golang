package ch8

import (
	"fmt"
	"log"
	"strings"

	"github.com/ianlopshire/go-fixedwidth"
)

type Book struct {
	ISBN        string
	PublishDate string
	Price       string
	PDF         string
	EPUB        string
	EbookPrice  string
}

var s string = `978-4-87311-865-9201909174620true true 3696
978-4-87311-924-3202010102750falsefalse0000
978-4-87311-878-9201903120000true true 0000`

func fixedLengthData() {
	fixedLengthData1()
	fixedLengthData2()
}

func fixedLengthData1() {
	for _, line := range strings.Split(s, "\n") {
		r := []rune(line)
		res := Book{
			ISBN:        string(r[0:17]),
			PublishDate: string(r[17:25]),
			Price:       string(r[25:29]),
			PDF:         string(r[29:34]),
			EPUB:        string(r[34:39]),
			EbookPrice:  string(r[39:43]),
		}
		fmt.Printf("%+v\n", res)
	}
}

type Book2 struct {
	ISBN        string `fixed:"1,17"`
	PublishDate string `fixed:"18,25"`
	Price       int    `fixed:"26,29"`
	PDF         string `fixed:"30,34,left"`
	EPUB        string `fixed:"35,39,left"`
	EbookPrice  int    `fixed:"40,44"`
}

func fixedLengthData2() {
	for _, line := range strings.Split(s, "\n") {
		var b Book2
		if err := fixedwidth.Unmarshal([]byte(line), &b); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", b)
	}
}
