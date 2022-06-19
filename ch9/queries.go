package ch9

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	UserID    string
	UserName  string
	CreatedAt time.Time
}

func queryMultiLines() {
	ctx := context.Background()
	cmd := `SELECT user_id, user_name, created_at FROM users ORDER BY user_id;`
	rows, err := db.QueryContext(ctx, cmd)
	if err != nil {
		log.Fatalf("query all users: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var (
			userID    string
			userName  string
			createdAt time.Time
		)
		if err := rows.Scan(&userID, &userName, &createdAt); err != nil {
			log.Fatalf("scan thet user; %v", err)
		}
		users = append(users, &User{
			UserID:    userID,
			UserName:  userName,
			CreatedAt: createdAt,
		})
	}
	if err := rows.Close(); err != nil {
		log.Fatalf("rows close: %v", err)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("scan users: %v", err)
	}

	for _, u := range users {
		fmt.Println(*u)
	}
}

func querySingleLine() {
	userID := `0001`
	var (
		userName  string
		createdAt sql.NullTime
	)
	ctx := context.Background()
	cmd := `SELECT user_name, created_at FROM users WHERE user_id = $1;`
	row := db.QueryRowContext(ctx, cmd, userID)
	err := row.Scan(&userName, &createdAt)
	if err != nil {
		log.Fatalf("query row(user_id=%s): %v", userID, err)
	}

	fmt.Println(createdAt.Valid)

	u := User{
		UserID:    userID,
		UserName:  userName,
		CreatedAt: createdAt.Time,
	}
	fmt.Println(u)
}
