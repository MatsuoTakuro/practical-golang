package ch1

import (
	"fmt"
	"io"
	"log"
	"os"
)

func defers() {
	err := closefileMyself()
	if err != nil {
		log.Fatalln(err)
	}

	err = deferReturnSample("/Users/user/training/go/practical-golang/ch1/defer.txt")
	if err != nil {
		log.Fatalln(err)
	}
}

var files []string

func closefileMyself() (err error) {
	files = append(files, "/Users/user/training/go/practical-golang/ch1/defer.go")
	var result [][]byte

	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("error when opening a file %w", err)
		}
		// defer f.Close()
		data, _ := io.ReadAll(f)
		result = append(result, data)
		f.Close()
	}
	for _, data := range result {
		fmt.Println(string(data))
	}
	return
}

func deferReturnSample(fname string) (err error) {
	var f *os.File
	f, err = os.Create(fname)
	if err != nil {
		return fmt.Errorf("error when opening a file %w", err)
	}
	defer func() {
		err = f.Close()
	}()
	_, err = io.WriteString(f, "A sample to catch defer's error")
	if err != nil {
		return fmt.Errorf("error when writing to a file %w", err)
	}
	return
}
