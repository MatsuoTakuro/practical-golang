package ch16

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
)

func workerPool() {
	TotalFileSize()
}

type Task string

type Result struct {
	value int64
	Task  Task
	Err   error
}

func TotalFileSize() int64 {
	tasks := make(chan Task)
	results := make(chan Result)

	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(i, tasks, results) // wait to receive a task and transmit a result
	}

	inputDone := make(chan struct{})
	var remainedCount int64
	go func() {
		path, _ := os.Getwd()
		err := filepath.Walk(string(path), func(path string, info fs.FileInfo, err error) error {
			atomic.AddInt64(&remainedCount, 1)
			tasks <- Task(path) // transmit a task
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
		close(inputDone)
		close(tasks)
	}()

	var totalSize int64
	for {
		select {
		case result := <-results: // wait to receive a result
			if result.Err != nil {
				fmt.Printf("err %v for %s\n", result.Err, result.Task)
			} else {
				atomic.AddInt64(&totalSize, result.value)
			}
			atomic.AddInt64(&remainedCount, -1)
		case <-inputDone: // wait to receive a inputDone
			if remainedCount == 0 {
				return totalSize
			}
		}
	}
}

func worker(id int, tasks <-chan Task, results chan<- Result) { // wait to receive a task
	for t := range tasks {
		fmt.Printf("worker: %d task: %s\n", id, t)

		fi, err := os.Stat(string(t))

		if err == nil && fi.IsDir() {
			err = fmt.Errorf("worker: %d err: %s is dir", id, fi)
		}

		result := Result{
			value: 0,
			Task:  t,
			Err:   nil,
		}

		if err != nil {
			result.Err = err
		} else {
			fmt.Printf("worker: %d path: %s size: %d\n", id, string(t), fi.Size())
			result.value = fi.Size()
		}

		results <- result // transmit a result
	}
}
