package ch16_test

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestLeak(t *testing.T) {
	wait := make(chan struct{})
	go func() {
		<-wait
	}()
}
