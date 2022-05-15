package ch3

import "fmt"

//go:generate enumer -type=Status

func structDesign() {
	ptrType()
	guaranteedZeroValue()
}

func ptrType() {
	n := NewNocopyStruct("test value")
	fmt.Println(n)
	fmt.Println(n.Value)
	n2 := n.Copy()
	fmt.Println(n2)
	fmt.Println(n2.Value)

	fmt.Println(n.Value == n2.Value)
	fmt.Println(*(n.Value) == *(n2.Value))
}

type NoCopyStruct struct {
	self  *NoCopyStruct
	Value *string
}

func NewNocopyStruct(value string) *NoCopyStruct {
	r := &NoCopyStruct{
		Value: &value,
	}
	r.self = r
	return r
}

func (n *NoCopyStruct) String() string {
	if n != n.self {
		panic("Should not create or copy NoCopyStruct instance without factory or Copy() method.")
	}
	return *n.Value
}

func (n *NoCopyStruct) Copy() *NoCopyStruct {
	if n != n.self {
		panic("Should not create or copy NoCopyStruct instance without factory or Copy() method.")
	}
	v := *n.Value
	return NewNocopyStruct(v)
}

type Status int

const (
	DefaultStatus Status = iota
	ActiveSatus
	CloseStatus
)

type visitor struct {
	Status Status
}

func guaranteedZeroValue() {
	var v visitor
	fmt.Println(v)
	fmt.Println(v.Status)
}
