package ch1

import (
	"log"
	"strings"
)

func strings_builder() {
	src := []string{"Golang", "Java", "Kotlin", "Python", "Ruby", "Rust"}
	var builder strings.Builder
	builder.Grow(100) // assumed that 100 or less chars should be used.
	for i, word := range src {
		if i != 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(word)
	}
	log.Println(builder.String())
}
