package ch16

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"

	"go.uber.org/ratelimit"
)

func workerPool() {

	// fmt.Println("TotalFileSize()", TotalFileSize())

	// taskSrcs := []Task{
	// 	"/Users/user/training/go/practical-golang/ch16/goroutine.go",
	// 	"/Users/user/training/go/practical-golang/ch16/goroutinepool.go",
	// 	"/Users/user/training/go/practical-golang/ch16/lock.go",
	// 	"/Users/user/training/go/practical-golang/ch16/sub.go",
	// }
	// fmt.Println("TotalFileSizeWithFixedTasks()", TotalFileSizeWithFixedTasks(taskSrcs))

	poolWithRateLimit()
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
		go fileSizeCalculator(i, tasks, results) // wait to receive a task and transmit a result
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

func TotalFileSizeWithFixedTasks(taskSrcs []Task) int64 {
	tasks := make(chan Task, len(taskSrcs))
	results := make(chan Result)
	for _, src := range taskSrcs {
		tasks <- src
	}
	close(tasks)

	for i := 0; i < runtime.NumCPU(); i++ {
		go fileSizeCalculator(i, tasks, results)
	}

	var count int
	var totalSize int64
	for {
		result := <-results
		count += 1
		if result.Err != nil {
			fmt.Printf("err %v for %s\n", result.Err, result.Task)
		} else {
			atomic.AddInt64(&totalSize, result.value)
		}

		if count == len(taskSrcs) {
			break
		}
	}
	return totalSize
}

func fileSizeCalculator(id int, tasks <-chan Task, results chan<- Result) { // wait to receive a task
	for t := range tasks {
		fmt.Printf("calculator: %d task: %s\n", id, t)

		fi, err := os.Stat(string(t))

		if err == nil && fi.IsDir() {
			err = fmt.Errorf("calculator: %d err: %s is dir", id, fi)
		}

		result := Result{
			value: 0,
			Task:  t,
			Err:   nil,
		}

		if err != nil {
			result.Err = err
		} else {
			fmt.Printf("calculator: %d path: %s size: %d\n", id, string(t), fi.Size())
			result.value = fi.Size()
		}

		results <- result // transmit a result
	}
}

func poolWithRateLimit() {
	tasks := make(chan Task, 100)
	results := make(chan Result)

	rl := ratelimit.New(100)
	for i := 0; i < runtime.NumCPU(); i++ {
		go workerWithRateLimit(rl, tasks, results)
	}

	inputDone := make(chan struct{})
	var remainedCount int64
	go func() {
		sampleTasks := []Task{"first", "second", "third"}

		for _, t := range sampleTasks {
			atomic.AddInt64(&remainedCount, 1)
			tasks <- t
		}
		close(inputDone)
		close(tasks)
	}()

	for {
		select {
		case result := <-results:
			if result.Err != nil {
				fmt.Printf("err %v for %s\n", result.Err, result.Task)
			} else {
				fmt.Printf("task: %s api request was successful(%d)!\n", result.Task, result.value)
			}
			atomic.AddInt64(&remainedCount, -1)
		case <-inputDone:
			if remainedCount == 0 {
				fmt.Printf("calling api is finised.\n")
				return
			}
		}
	}
}

func workerWithRateLimit(rt ratelimit.Limiter, tasks <-chan Task, results chan<- Result) {
	for t := range tasks {
		fmt.Printf("Let's get started to work on task: %s!\n", t)

		rt.Take()
		v, err := callExternalAPI()
		if err != nil {
			err = fmt.Errorf("task: %s error: %v", t, err)
		}

		result := Result{
			value: v,
			Task:  t,
			Err:   err,
		}
		results <- result
	}
}

func callExternalAPI() (value int64, err error) {
	resp, err := http.Get("https://example.com")

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var v int64
	if v = int64(resp.StatusCode); v != 200 {
		return 0, err
	}

	return v, nil
}
