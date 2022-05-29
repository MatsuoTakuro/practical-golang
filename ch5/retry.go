package ch5

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/avast/retry-go"
)

func random() error {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	println("random:", n)
	if n >= 5 {
		return nil
	}
	return errors.New("error")
}

func withRetry() {
	err := retry.Do(
		// retry
		func() error {
			ret := random()
			return ret
		},
		// opt1 (default is BackOff)
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			fmt.Println("delay(sec):", n)
			return time.Duration(n) * time.Second
		}),
		// opt2 (default is 10)
		retry.Attempts(3),
	)
	fmt.Println(err)
}
