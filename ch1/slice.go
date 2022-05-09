package ch1

import "fmt"

func slice() {
	s1 := make([]int, 1000)
	fmt.Println(len(s1))
	fmt.Println(cap(s1))

	s2 := make([]int, 0, 1000)
	fmt.Println(len(s2))
	fmt.Println(cap(s2))

}
