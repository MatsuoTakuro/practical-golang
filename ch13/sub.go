package ch13

import "fmt"

func Sub() {
	// server()
	fmt.Println("AppendSlice", AppendSlice(5, 1))
	fmt.Println("FirstAllocSlice", FirstAllocSlice(5, 1))
}
