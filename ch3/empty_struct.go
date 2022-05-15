package ch3

import "fmt"

func emptyStruct() {
	wait := make(chan struct{})
	go func() {
		fmt.Println("sent")
		wait <- struct{}{}
	}()

	fmt.Println("waiting for receipt")
	<-wait
	fmt.Println("received")

}
