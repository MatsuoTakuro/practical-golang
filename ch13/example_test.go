package ch13

import "fmt"

// godoc -http=localhost:6060
func ExampleF() {
	F()
	// output: example_f
}

func ExampleT() {
	t := T{}
	fmt.Println(t)
	// output: example_t
}

func ExampleT_M() {
	t := T{}
	t.M()
	// output: example_m
}
