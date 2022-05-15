package ch2

import (
	"container/list"
	"fmt"
	"net/url"
)

func with_method() {
	httpStatus()
	urlValues()
	containerList()
}

func httpStatus() {
	status := StatusOK
	fmt.Println(status.String())
}

type HTTPStatus int

const (
	StatusOK              HTTPStatus = 200
	StatusUnauthorized    HTTPStatus = 401
	StatusPaymentRequired HTTPStatus = 402
	StatusForbidden       HTTPStatus = 403
)

func (s HTTPStatus) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusPaymentRequired:
		return "Payment Required"
	case StatusForbidden:
		return "Forbidden"
	default:
		return fmt.Sprintf("HTTPStatus(%d)", s)
	}
}

func urlValues() {
	vs := url.Values{}
	vs.Add("key1", "value1")
	vs.Add("key2", "value2")
	for k, v := range vs {
		fmt.Printf("%s: %v\n", k, v)
	}
}

// like iterator
func containerList() {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	for ele := l.Front(); ele != nil; ele.Next() {
		fmt.Println(ele.Value)
	}
}
