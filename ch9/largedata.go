package ch9

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/jackc/pgx/v4"
)

func preparedStmt() {
	users := []User{
		{
			UserID:   "0003",
			UserName: "Duke",
		},
		{
			UserID:   "0004",
			UserName: "John",
		},
		{
			UserID:   "0005",
			UserName: "Micheal",
		},
	}

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("begin trans: %v\n", err)
	}
	defer tx.Rollback()

	cmd := `INSERT INTO users(user_id, user_name, created_at)VALUES( $1, $2, $3 );`
	stmt, err := tx.PrepareContext(ctx, cmd)
	if err != nil {
		log.Fatalf("create prepared stmt: %v\n", err)
	}
	defer stmt.Close()

	for _, u := range users {
		if _, err := stmt.ExecContext(ctx, u.UserID, u.UserName, nil); err != nil {
			log.Fatalf("execute stmt: %v\n", err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("commit trans: %v\n", err)
	}
}

func batchInsert() {
	users := []User{
		{
			UserID:   "0006",
			UserName: "Luke",
		},
		{
			UserID:   "0007",
			UserName: "Misa",
		},
		{
			UserID:   "0008",
			UserName: "Dior",
		},
	}

	ctx := context.Background()

	valueStrings := make([]string, 0, len(users)) // 3
	valueArgs := make([]any, 0, len(users)*2)     // 3 * 2 = 6

	number := 1
	for _, u := range users {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", number, number+1)) // ($1, $2) -> ($3, $4) -> ($5, $6)
		valueArgs = append(valueArgs, u.UserID)
		valueArgs = append(valueArgs, u.UserName)
		number += 2
	}

	cmd := fmt.Sprintf("INSERT INTO users (user_id, user_name) VALUES %s;", strings.Join(valueStrings, ","))
	fmt.Println(cmd)
	fmt.Println(valueArgs)
	if _, err := db.ExecContext(ctx, cmd, valueArgs...); err != nil {
		log.Fatalf("execute query: %v\n", err)
	}
}

func builtInDbFuncs() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	conn, err := pgx.Connect(context.Background(), "postgres://testuser:pass@localhost:5432/testdb")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rows := [][]any{
		{3, "おにぎり", 120},
		{4, "パン", 300},
		{5, "お茶", 100},
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"products"},
		[]string{"product_no", "name", "price"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
