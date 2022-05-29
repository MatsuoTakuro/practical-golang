package ch5

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/multierr"
)

func aggregateErrors() {
	// _ = allOrNothing()
	err := individualHandling()
	if err != nil {
		fmt.Println(err)
	}
}

type errWriter struct {
	w   *bufio.Writer
	err error
}

func (ew *errWriter) write(buf []byte) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}
func (ew *errWriter) flush() {
	if ew.err != nil {
		return
	}
	err := ew.w.Flush()
	if err != nil {
		ew.err = err
	}
}

func allOrNothing() error {
	f, _ := os.OpenFile("ch5/test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	writer := bufio.NewWriter(f)
	ss := []string{"falcon", "falcon2", "falcon3"}
	ew := &errWriter{
		w:   writer,
		err: nil,
	}
	for _, s := range ss {
		b := []byte(fmt.Sprintf("%s\n", s))
		ew.write(b)
	}
	if ew.err != nil {
		return ew.err
	}
	ew.flush()
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func individualHandling() error {
	var merr error
	for i := 0; i < 3; i++ {
		var err error
		if i%2 == 0 {
			_, err = os.Open(fmt.Sprintf("non-existing[%d]", i))
		}
		if err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	if merr != nil {
		return merr
	}
	return nil
}
