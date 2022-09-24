package ch16

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func lock() {
	lockAndUnlock()
	atomicInteger()
}

func lockAndUnlock() {
	ac := &Account{
		balance: 0,
		lock:    sync.RWMutex{},
	}
	fmt.Println("initialized balance:", ac.GetBalance())

	ic := make(chan int)
	go func() {
		ic <- 100 // transmit
		ic <- 300 // transmit
		close(ic)
	}()

	ic2 := make(chan int)
	go func() {
		ic2 <- 200 // transmit
		ic2 <- 400 // transmit
		close(ic2)
	}()

	for am := range ic2 { // receive * n
		ac.Transfer(am)
	}
	fmt.Println("transferred balance:", ac.GetBalance())

	for am := range ic { // receive * n
		ac.Transfer(am)
	}
	fmt.Println("transferred balance:", ac.GetBalance())

}

type Account struct {
	balance int
	lock    sync.RWMutex
}

func (a *Account) GetBalance() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.balance
}

func (a *Account) Transfer(amount int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.balance += amount
	time.Sleep(1 * time.Millisecond)
}

func atomicInteger() {
	ac := &Account3{
		balance: 0,
	}
	fmt.Println("initialized balance:", ac.GetBalance())

	i64c := make(chan int64)
	go func() {
		i64c <- 100 // transmit
		i64c <- 200 // transmit
		close(i64c)
	}()

	i64c2 := make(chan int64)
	go func() {
		i64c2 <- 100 // transmit
		i64c2 <- 200 // transmit
		close(i64c2)
	}()

	for am := range i64c2 { // receive * n
		ac.Transfer(am)
	}
	fmt.Println("transferred balance:", ac.GetBalance())

	for am := range i64c { // receive * n
		ac.Transfer(am)
	}
	fmt.Println("transferred balance:", ac.GetBalance())
}

type Account3 struct {
	balance int64
}

func (a Account3) GetBalance() int64 {
	return a.balance
}

func (a *Account3) Transfer(amount int64) {
	atomic.AddInt64(&a.balance, amount)
	time.Sleep(1 * time.Millisecond)
}
