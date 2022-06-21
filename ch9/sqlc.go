package ch9

import (
	"context"
	"fmt"
	"log"

	"github.com/MatsuoTakuro/practical-golang/ch9/testdb"
)

func sqlc() {
	ctx := context.Background()
	queries := testdb.New(db)

	us, err := queries.ListUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range us {
		fmt.Printf("%+v\n", u)
	}
}
