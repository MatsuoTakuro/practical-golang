package ch4

import (
	"fmt"

	"go.uber.org/dig"
)

var Container *dig.Container

type RandomGenerator interface {
	Generate() float64
}

type fakeRandomGenerator struct{}

func (f fakeRandomGenerator) Generate() float64 {
	return 1.0
}

func init() {
	Container = dig.New()

	Container.Provide(func() RandomGenerator {
		return &fakeRandomGenerator{}
	})
}

func uberDig() {
	Container.Invoke(func(rg RandomGenerator) {
		fmt.Println("random number:", rg.Generate())
	})
}
