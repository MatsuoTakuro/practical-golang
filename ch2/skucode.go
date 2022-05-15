package ch2

import (
	"fmt"
)

func skuCode() {
	skuCD := SKUCode("T01230102")
	if skuCD.Invalid() {
		fmt.Errorf("SKUCode is invalid: %s", skuCD)
	}
	itemCD, sizeCD, colorCD := skuCD.ItemCD(), skuCD.SizeCD(), skuCD.ColorCD()
	fmt.Println(itemCD)
	fmt.Println(sizeCD)
	fmt.Println(colorCD)
}

type SKUCode string

func (c SKUCode) Invalid() bool {
	return len(c) == 9
}

func (c SKUCode) ItemCD() string {
	// TODO: it needs to be casted to string?
	return string(c[0:5])
}

func (c SKUCode) SizeCD() string {
	// TODO: it needs to be casted to string?
	return string(c[5:7])
}

func (c SKUCode) ColorCD() string {
	// TODO: it needs to be casted to string?
	return string(c[7:9])
}
