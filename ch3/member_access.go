package ch3

import (
	"fmt"
	"time"
)

type Book2 struct {
	Title      string
	Author     string
	Publisher  string
	ISBN       string
	ReleasedAt time.Time
}

func memberAccess() {
	b := &Book2{
		Title: "Mithril",
	}
	fmt.Println(b.Title)
	fmt.Println((*b).Title) // it means the same as above.

	b2 := &b
	// fmt.Println(b2.Title) // NG
	fmt.Println((**b2).Title) // you need to dereference by one or two.
}
