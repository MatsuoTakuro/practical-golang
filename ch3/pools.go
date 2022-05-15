package ch3

import (
	"fmt"
	"sync"
)

func pools() {
	b := pool.Get().(*BigStruct)
	fmt.Println(b.Member)
	b.Member = "Taro"
	fmt.Println(b.Member)
	pool.Put(b)
	rb := pool.Get().(*BigStruct)
	fmt.Println(rb.Member)

	b2 := NewBigStruct2()
	fmt.Println(b2.Member)
	b2.Member = "Taro2"
	fmt.Println(b2.Member)
	b2.Release()
	rb2 := NewBigStruct2()
	fmt.Println(rb2.Member)
}

type BigStruct struct {
	Member string
}

var pool = &sync.Pool{
	New: func() any {
		return &BigStruct{}
	},
}

type BigStruct2 struct {
	Member string
}

var pool2 = &sync.Pool{
	New: func() any {
		return &BigStruct2{}
	},
}

func NewBigStruct2() *BigStruct2 {
	b2 := pool2.Get().(*BigStruct2)
	return b2
}

func (b2 *BigStruct2) Release() {
	b2.Member = ""
	pool2.Put(b2)
}
