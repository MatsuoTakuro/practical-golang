package ch16

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func fanOutAndFanIn() {
	// syncWaitGroup()
	errgroupGroup()
}

func syncWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Done: 1")
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Done: 2")
		wg.Done()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Done: 3")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Done all tasks!")
}

func errgroupGroup() {
	eg, ctx := NewErrorGroup()

	eg.Go(func() error {
		cmd := exec.CommandContext(ctx, "sleep", "1")
		err := cmd.Run()
		fmt.Println("Done: 1")
		return err
	})

	eg.Go(func() error {
		cmd := exec.CommandContext(ctx, "sleep", "2")
		err := cmd.Run()
		fmt.Println("Done: 2")
		return err
	})

	eg.Go(func() error {
		cmd := exec.CommandContext(ctx, "sleep", "3")
		err := cmd.Run()
		fmt.Println("Done: 3")
		return err
	})

	eg.Go(func() error {
		time.Sleep(2 * time.Second)
		return errors.New("accidental error")
	})

	err := eg.Wait()
	fmt.Printf("result=> err: %v, context: %v\n", err, ctx)
}

func NewErrorGroup() (*errgroup.Group, context.Context) {
	return errgroup.WithContext(context.Background())
}
