package ch10

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type Book struct {
	Title string `validate:"required"`
	Price *int   `validate:"required"`
}

func (b Book) String() string {
	return fmt.Sprintf("Title:%s Price:%d\n", b.Title, *b.Price)
}

func required() {
	s := `{"Title":"Real World HTTP ミニ版", "Price":0}` // 無料版のためPriceは0円
	var b Book
	if err := json.Unmarshal([]byte(s), &b); err != nil {
		log.Fatal(err)
	}

	if err := validator.New().Struct(b); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fmt.Printf("field %s violates %s (value:%v)", fe.Field(), fe.Tag(), fe.Value())
			}
		}
	}
	fmt.Printf("no violation on %+v", b)
}
