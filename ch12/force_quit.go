package ch12

import (
	"fmt"
	"log"
)

func forceQuit() {
	logFatal()
	// panicAndRecover()
}

func logFatal() {
	fmt.Println(1)
	defer func() {
		fmt.Println(3)
		err := recover()
		fmt.Println("in-recover:", err)
		fmt.Println(4)
	}()

	fmt.Println(2)
	log.Fatal("fatal message")
}

func panicAndRecover() {
	fmt.Println(1)
	defer func() {
		fmt.Println(3)
		err := recover()
		fmt.Println("in-recover:", err)
		fmt.Println(4)
	}()

	fmt.Println(2)
	panic("panic message")
}
