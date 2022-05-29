package ch5

import (
	"fmt"

	"golang.org/x/xerrors"
)

func stackTrace() {

	err := xerrors.New("xerrors error")
	fmt.Printf("%+v\n", err)
}
