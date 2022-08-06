package ch13

// go test -benchmem -bench (AppendSlice|BenchmarkFirstAllocSlice) github.com/MatsuoTakuro/practical-golang/ch13
func AppendSlice(count, value int) []int {
	res := []int{}
	for i := 0; i < count; i++ {
		res = append(res, value)
	}
	return res
}

func FirstAllocSlice(count, value int) []int {
	res := make([]int, count)
	for i := 0; i < count; i++ {
		// res = append(res, value)
		res[i] = value
	}
	return res
}
