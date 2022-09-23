package ch16

import (
	"fmt"
	"time"
)

func goroutine() {
	loop1()
	loop2()
}

func loop1() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, v := range items {
		v2 := v
		go func() {
			fmt.Printf("v2 = %d, address = %p\n", v2, &v2)
		}()
	}
	time.Sleep(time.Second)
}

func loop2() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, v := range items {
		go func(v int) {
			fmt.Printf("v = %d, address = %p\n", v, &v)
		}(v)
	}
	time.Sleep(time.Second)
}
