package ch16

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

func errorEventWithContext() {
	ctx := context.Background()
	err := runJobsWithError(ctx)
	fmt.Println("runJobs err:", err)
}

func runJobsWithError(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // when error happens or all tasks are done
	ec := make(chan error)
	done := make(chan struct{})

	for i := 0; i < 10; i++ {
		go func() {
			cmd := exec.CommandContext(ctx, "sleep", "3")
			err := cmd.Run()
			if err != nil {
				ec <- err // transmit
			} else {
				done <- struct{}{} // transmit
			}
		}()
	}

	go func() {
		time.Sleep(1 * time.Second)
		ec <- errors.New("accidental error") // transmit
	}()

	for i := 0; i < 11; i++ {
		select {
		case err := <-ec: // receive
			return err
		case <-done: // receive
			fmt.Printf("cmd #%d is done.\n", i)
		}
	}
	return nil
}
