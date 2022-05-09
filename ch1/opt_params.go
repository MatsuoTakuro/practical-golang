package ch1

import (
	"fmt"
	"time"
)

type Portion int

const (
	Regular Portion = iota + 1
	Small
	Large
)

type Udon struct {
	men      Portion
	aburaage bool
	ebiten   uint
}

func opt_params() {
	opt := &Option{Large, true, 2}
	udon := newUdonByStruct(*opt)
	fmt.Println(*udon)

	udon2 := newUdonByBuilder(Large).Aburaage().Ebiten(2).Order()
	fmt.Println(*udon2)

	udon3 := newUdonByFunctionalBuilder(OptMen(Large), OptAburaage(), OptEbiten(2))
	fmt.Println(*udon3)
}

type Option struct {
	men      Portion
	aburaage bool
	ebiten   uint
}

func newUdonByStruct(opt Option) *Udon {
	if opt.ebiten == 0 && time.Now().Hour() < 10 { // the range of Time.Hour() is [0, 23].
		opt.ebiten = 1
	}
	return &Udon{
		men:      opt.men,
		aburaage: opt.aburaage,
		ebiten:   opt.ebiten,
	}
}

type fluentOpt struct {
	men      Portion
	aburaage bool
	ebiten   uint
}

func newUdonByBuilder(p Portion) *fluentOpt { // Portion is a required option.
	return &fluentOpt{
		men:      p,
		aburaage: false,
		ebiten:   1,
	}
}

func (o *fluentOpt) Aburaage() *fluentOpt {
	o.aburaage = true
	return o
}

func (o *fluentOpt) Ebiten(n uint) *fluentOpt {
	o.ebiten = n
	return o
}

func (o *fluentOpt) Order() *Udon {
	return &Udon{
		men:      o.men,
		aburaage: o.aburaage,
		ebiten:   o.ebiten,
	}
}

type OptFunc func(u *Udon)

func newUdonByFunctionalBuilder(opts ...OptFunc) *Udon {
	u := &Udon{}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func OptMen(p Portion) OptFunc {
	return func(u *Udon) { u.men = p }
}

func OptAburaage() OptFunc {
	return func(u *Udon) { u.aburaage = true }
}

func OptEbiten(n uint) OptFunc {
	return func(u *Udon) { u.ebiten = n }
}
