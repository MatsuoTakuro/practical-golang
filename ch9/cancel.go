package ch9

import (
	"context"
	"log"
	"time"
)

func cancel() {
	s := Service{
		db: db,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, "SELECT pg_sleep(100);"); err != nil {
		log.Println("canceling query")
	} else {
		log.Fatalf("sleep: %v", err)
	}
}
