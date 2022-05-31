package ch6

import (
	"fmt"

	_ "github.com/MatsuoTakuro/practical-golang/ch6/b"
	_ "github.com/MatsuoTakuro/practical-golang/ch6/c"
)

func init() {
	fmt.Println("ch6.init")
}

func Sub() {
	fmt.Println("ch6")
}
