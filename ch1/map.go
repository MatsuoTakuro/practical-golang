package ch1

import "fmt"

func maps() {
	m := make(map[string]int, 1000)
	fmt.Println(len(m))
	// fmt.Println(cap(m)) // compile error

}
