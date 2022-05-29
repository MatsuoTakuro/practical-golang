package ch5

import (
	"errors"
	"log"
)

func errcheck() {
	_ = a()
	_, err := a2()
	if err != nil {
		log.Fatal(err)
	}
	_ = b()
	if err != nil {
		log.Fatal(err)
	}
}

func a() error {
	return errors.New("some error")
}

func a2() (string, error) {
	return "", nil
}

func b() error {
	return errors.New("b error")
}
