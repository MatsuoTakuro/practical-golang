package ch13

import (
	"testing"
)

func Add(a, b int) (int, error) {
	return a + b, nil
}

func TestAdd(t *testing.T) {
	got, _ := Add(1, 2)
	if got != 3 {
		t.Errorf("expect 3, but %d", got)
	}
}
