package ch16

import (
	"context"
	"fmt"
	"time"
)

func goroutineLeak() {
	// receiveInForStatement()
	// transmitInForStatement()
	closeChannelByCloseMethod()
}

func receiveInForStatement() {
	tasks := make(chan Task)

	go func() {
		tasks <- Task("xxx")
		close(tasks)
	}()

	go func() {
		defer fmt.Println("for-statement done")
		// NG
		// for {
		// 	task := <-tasks // forever wait to receive
		// 	fmt.Println(task)
		// }

		// OK
		for task := range tasks { // exit when chan is closed.
			fmt.Println(task)
		}
	}()

	time.Sleep(1 * time.Millisecond)
}

func transmitInForStatement() {
	ic := make(chan int)

	// NG
	// go func() {
	// 	defer fmt.Println("for-statement done")
	// 	count := 0
	// 	for {
	// 		ic <- count // transmit
	// 		count++
	// 	}
	// }()

	// OK
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		count := 0
		for {
			select {
			case ic <- count: // transmitted?
				ic <- count // transmit (again)
				count++
				fmt.Println("count", count)
			case <-ctx.Done(): // cancelled?
				fmt.Println("ctx.Done")
				return // timeout etc.
			}
		}
	}()

	go func() {
		for i := range ic { // receive
			fmt.Println("i", i)
		}
	}()

	time.Sleep(1 * time.Millisecond)
}

func closeChannelByCloseMethod() {
	ic := NewInfiniteCounter()
	fmt.Println(<-ic.Counter)
	fmt.Println(<-ic.Counter)
	fmt.Println(<-ic.Counter)
	ic.Close()
	fmt.Println(<-ic.Counter) // zero value cuz Continue channel is closed.
}

type InfiniteCounter struct {
	Counter chan int
	exit    chan struct{}
}

func NewInfiniteCounter() *InfiniteCounter {
	ic := &InfiniteCounter{
		Counter: make(chan int),
		exit:    make(chan struct{}),
	}

	go func() {
		count := 0
		for {
			select {
			case ic.Counter <- count:
				count++
			case <-ic.exit:
				close(ic.Counter)
				return
			}
		}
	}()
	return ic // result
}

func (ic *InfiniteCounter) Close() {
	close(ic.exit)
}
