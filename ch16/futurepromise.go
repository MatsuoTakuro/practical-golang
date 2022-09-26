package ch16

import (
	"log"
	"time"
)

func A() int {
	time.Sleep(time.Second)
	return 10
}

func B() int {
	time.Sleep(time.Second * 2)
	return 5
}

func C(a, b int) int {
	time.Sleep(time.Second * 1)
	return a + b
}

func D(a, c int) int {
	time.Sleep(time.Second)
	return a + c
}

func futureAndPromise() {
	fa, pa := MakeFuturePromise()
	fb, pb := MakeFuturePromise()
	fc, pc := MakeFuturePromise()
	fd, pd := MakeFuturePromise()

	go func() {
		a := A()
		pa.Submit(a)
	}()
	// time.Sleep(2 * time.Second) // wait for chan to be closed by Promise.Submit(v)
	log.Printf("a: %d value: %d done?: %t\n", fa.Get(), fa.value, fa.IsDone())

	go func() {
		b := B()
		pb.Submit(b)
	}()
	// time.Sleep(2 * time.Second) // wait for chan to be closed by Promise.Submit(v)
	log.Printf("b: %d value: %d done?: %t\n", fb.Get(), fb.value, fb.IsDone())

	go func() {
		c := C(fa.Get(), fb.Get())
		pc.Submit(c)
	}()
	// time.Sleep(2 * time.Second) // wait for chan to be closed by Promise.Submit(v)
	log.Printf("c: %d value: %d done?: %t\n", fc.Get(), fc.value, fc.IsDone())

	go func() {
		d := D(fa.Get(), fc.Get())
		pd.Submit(d)
	}()
	// time.Sleep(2 * time.Second) // wait for chan to be closed by Promise.Submit(v)
	log.Printf("d: %d value: %d done?: %t\n", fd.Get(), fd.value, fd.IsDone())
}

type Future struct {
	value int // input
	wait  chan struct{}
}

func (f *Future) IsDone() bool {
	select {
	case <-f.wait: // transmitted or closed channel?
		return true
	default:
		return false
	}
}

func (f *Future) Get() int {
	<-f.wait // transmit
	return f.value
}

type Promise struct {
	f *Future
}

func (p *Promise) Submit(v int) {
	p.f.value = v // set value as input of Future
	close(p.f.wait)
}

func MakeFuturePromise() (*Future, *Promise) {
	f := &Future{
		value: 0,
		wait:  make(chan struct{}),
	}
	p := &Promise{
		f: f,
	}
	return f, p
}
