package ch4

import "fmt"

type Fish any

func castedSlice() {
	var fishList = []Fish{"mackerel", "Japanese amberjack", "tunny"}
	// var fishNameList = fishList.([]string) // cannot directly downcast fishList([]Fish) to []string
	// var anyList []any = fishList           // cannot directly upcast fishList([]Fish) to []any

	// downcasting
	fishNames := make([]string, len(fishList))
	for i, f := range fishList {
		// type assertion
		if fn, ok := f.(string); ok {
			fishNames[i] = fn
		}
	}
	fmt.Println(fishList)

	// upcasting
	anyValues := make([]any, len(fishNames))
	for i, fn := range fishNames {
		// you donot need to do type assertion, but to assign the value to each element
		anyValues[i] = fn
	}
	fmt.Println(anyValues)
}
