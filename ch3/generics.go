package ch3

import (
	"fmt"
	"strconv"
)

func generics() {
	// https://zenn.dev/nobishii/articles/type_param_intro
	myInt := []MyInt{1, 2, 3, 4}
	fmt.Println(f(myInt))

	stack := NewStack[string]()
	stack.Push("hello")
	stack.Push("world")
	fmt.Println(stack.Pop()) // world
	fmt.Println(stack.Pop()) // hello

	set := NewSet(1, 2, 3)
	set.Add(5)
	fmt.Println(set)
	fmt.Println(set.Includes(3)) // true
	set.Remove(3)
	fmt.Println(set)
	fmt.Println(set.Includes(3)) // false

	// in this book
	var some Some
	some.i = "test"
	fmt.Println(some.String())

	fmt.Println(String(myInt))
	fmt.Println(String(stack))
	fmt.Println(String(set))
}

// f is a function with a Type Parameter.
// T is a Type Parameter.
// Stringer Inteface is used as Type Constraint against T.
func f[T Stringer](ss []T) []string {
	var result []string
	for _, s := range ss {
		// s can use String() method by implementing Stringer Inteface(used as Type Constraint).
		result = append(result, s.String())
	}
	return result
}

type Stringer interface {
	String() string
}

type MyInt int

// MyInt implements Stringer Inteface.
func (mi MyInt) String() string {
	return strconv.Itoa(int(mi))
}

// any means interface{}. So, this means []interface{}.
type Stack[T any] []T

func NewStack[T any]() *Stack[T] {
	v := make(Stack[T], 0)
	return &v
}

func (s *Stack[T]) Push(x T) {
	*s = append(*s, x)
}

func (s *Stack[T]) Pop() T {
	v := (*s)[len(*s)-1]  // pop the last ele.
	*s = (*s)[:len(*s)-1] // Substitute the slice up to the ele before the last one(= v).
	return v
}

// https://pkg.go.dev/builtin#comparable
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](cs ...T) Set[T] {
	s := make(Set[T])
	for _, c := range cs {
		s.Add(c)
	}
	return s
}

func (s Set[T]) Add(c T) {
	s[c] = struct{}{}
}

func (s Set[T]) Includes(c T) bool {
	_, ok := s[c]
	return ok
}

func (s Set[T]) Remove(c T) {
	delete(s, c)
}

// String(x T) is a function with a Type Parameter.
func String[T any](x T) string {
	return fmt.Sprintf("%v by T", x)
}

type Some struct {
	i interface{}
}

func (s Some) String() string {
	return fmt.Sprintf("%s by Some", String(s.i))
}

// type Some2[T any] struct {
// 	t T
// }

// method must have no type parameters
// func (s2 Some2[T]) Method[R any](r R) {
// 	fmt.Println(s2.t, r)
// }

// type Some3 struct {
// 	name string
// }

// method must have no type parameters
// func (s3 Some3) Method[T any](x T) {
// 	fmt.Println(s3, x)
// }
