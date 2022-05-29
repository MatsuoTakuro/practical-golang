package ch5

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Refer to https://zenn.dev/msksgm/articles/20220325-unwrap-errors-is-as
func unwrap() {
	basic()
	withErrorf()
	withOwnDefinedError()
}

func testErr() error {
	return errors.New("test error")
}
func basic() {
	err := fmt.Errorf("wrapping: %w", testErr())
	fmt.Println(err)
	fmt.Println(errors.Unwrap(err))
	fmt.Println(testErr())
}

func withErrorf() {
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *fs.PathError

		wrapedErr := fmt.Errorf("err is %w", err)
		if errors.As(wrapedErr, &pathError) {
			fmt.Println("errors.As():Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
		if errors.Is(wrapedErr, pathError) {
			fmt.Println("errors.Is():Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}
}

type SampleError struct {
	msg string
	err error
}

func (se *SampleError) Error() string {
	return se.msg
}
func (se *SampleError) Unwrap() error { // Need to add
	return se.err
}
func withOwnDefinedError() {
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *fs.PathError

		wrapedErr := &SampleError{msg: "this is wraped err", err: err}
		if errors.As(wrapedErr, &pathError) {
			fmt.Println("errors.As():Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
		if errors.Is(wrapedErr, pathError) {
			fmt.Println("errors.Is():Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}
}
