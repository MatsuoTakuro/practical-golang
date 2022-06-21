package ch9

import (
	"context"
	"fmt"
	"log"

	"github.com/MatsuoTakuro/practical-golang/ch9/models"
)

func sqlboiler() {
	ctx := context.Background()
	userID := "0001"

	where(ctx, userID)
	find(ctx, userID)
	all(ctx)
}

func where(ctx context.Context, userID string) {
	u, err := models.Users(models.UserWhere.UserID.EQ(userID)).One(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", *u)
}

func find(ctx context.Context, userID string) {
	u, err := models.FindUser(ctx, db, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", *u)
}

func all(ctx context.Context) {
	us, err := models.Users().All(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range us {
		fmt.Printf("%+v\n", *u)
	}
}
