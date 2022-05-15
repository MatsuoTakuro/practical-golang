package ch2

import (
	"fmt"
	"sort"
	"time"
)

func consumers() {
	cs := NewSamleConsumers()
	fmt.Println(cs)
	fmt.Println(cs.RequiredFollows())
	fmt.Println(cs)

}

type Consumer struct {
	Name       string
	ActiveFlag bool
	ExpiredAt  time.Time
}

type Consumers []Consumer

func NewSamleConsumers() Consumers {
	return Consumers{
		Consumer{"John", true, time.Now().AddDate(0, 1, 0)},
		Consumer{"Smith", true, time.Now().AddDate(0, 5, 0)},
		Consumer{"Max", true, time.Now().AddDate(0, 3, 0)},
	}
}

func (cs Consumers) RequiredFollows() Consumers {
	return cs.activeConsumer().expires(time.Now().AddDate(0, 2, 0)).sortByExpiredAt()
}

func (cs Consumers) activeConsumer() Consumers {
	resp := make([]Consumer, 0, len(cs))
	for _, c := range cs {
		if c.ActiveFlag {
			resp = append(resp, c)
		}
	}
	return resp
}

func (cs Consumers) expires(end time.Time) Consumers {
	resp := make([]Consumer, 0, len(cs))
	for _, c := range cs {
		if c.ExpiredAt.After(end) {
			resp = append(resp, c)
		}
	}
	return resp
}

func (cs Consumers) sortByExpiredAt() Consumers {
	sort.Slice(cs, func(c1, c2 int) bool {
		return cs[c1].ExpiredAt.Before(cs[c2].ExpiredAt)
	})
	return cs
}
