package ch4

import (
	"context"
	"log"
)

func contexts() {
	ctx := context.WithValue(context.Background(), "favorite", "Zenigata Heiji")

	if s, ok := ctx.Value("favorite").(string); ok {
		log.Printf("My favorite is %s.", s)
	}

	switch v := ctx.Value("favorite").(type) {
	case string:
		log.Printf("My favorite thing: %s\n", v)
	case int:
		log.Printf("My favorite number: %d\n", v)
	case complex128:
		log.Printf("My favorite complex number: %f\n", v)
	default:
		log.Printf("My favorite: %v\n", v)

	}
}
