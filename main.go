package main

import "fmt"

//go:generate stringer -type=CarType
//go:generate stringer -type=CarOption

type CarType int
type CarOption uint64

const (
	Sedan CarType = iota + 1
	Hatchback
	MPV
	SUV
	Crossover
	Coupe
	Convertible
)

const (
	GPS          CarOption = 1 << iota // 01(2)    = 1(10)
	AWD                                // 10(2)    = 2(10)
	SunRoof                            // 100(2)   = 4(10)
	HeatedSeat                         // 1000(2)  = 8(10)
	DriverAssist                       // 10000(2) = 16(10)
)

func main() {
	carType()
	carOption()
}

func carType() {
	var t CarType
	fmt.Println(t)
	t = Sedan
	fmt.Println(t)
}

func carOption() {
	var o CarOption
	fmt.Println(o)
	o = SunRoof | HeatedSeat // 100(2)  | 1000(2) = 1100(2) = 12(10)
	if o&SunRoof != 0 {      // 1100(2) & 1000(2) =  100(2) = 4(10)
		fmt.Printf("with %s\n", o)
	}
}
