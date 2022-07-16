package ch12

import (
	"io"
	"log"
	"os"
)

func logOutput() {
	// standardLib()
	// multiOutput()
	customization()
}

func standardLib() {
	log.Println("Log output")
	n := 10
	s := "string"
	c := 1 + 2i
	log.Printf("variables output can also be done using %d, %s. %v can show any type.", n, s, c)
}

func multiOutput() {
	file, _ := os.Create("ch12/log.txt")
	log.SetOutput(io.MultiWriter(file, os.Stderr))
	log.Println("Log outputs both file and stderr at the same time")
}

func customization() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("🐙 ")
	log.Println("Log output")
}
