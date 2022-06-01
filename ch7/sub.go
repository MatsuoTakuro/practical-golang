package ch7

import (
	"fmt"
	"runtime/debug"
)

func Sub(version string) {
	fmt.Printf("version: %s\n", version)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	str := info.String()
	fmt.Println(str)

	panic("panic in here")
}
