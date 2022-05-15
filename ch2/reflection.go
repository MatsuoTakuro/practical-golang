package ch2

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func reflection() {
	writeType := reflect.TypeOf((*io.Writer)(nil)).Elem()
	writeType2 := reflect.TypeOf((*io.Writer)(nil)).Elem()
	fileType := reflect.TypeOf((*os.File)(nil)).Elem()
	fmt.Println(fileType.Implements(writeType))
	fmt.Println(writeType.Implements(writeType2))
}
