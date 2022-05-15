package ch2

func season() {
	p := Peak
	println(int(p.Price(100)))
	n := Normal
	println(int(n.Price(100)))

}

type Season int

const (
	Peak Season = iota + 1
	Normal
	Off
)

func (s Season) Price(price float64) float64 {
	if s == Peak {
		return price + 200
	}
	return price
}
