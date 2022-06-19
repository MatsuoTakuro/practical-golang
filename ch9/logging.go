package ch9

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/qustavo/sqlhooks/v2"
)

// implement Logger interface
type logger struct{}

func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if msg == "Query" {
		log.Printf("SQL:\n%v\nARGS:%v\n", data["sql"], data["args"])
	}
}

var _ pgx.Logger = (*logger)(nil)

type PgTable struct {
	SchemaName string
	TableName  string
}

func loggingWithDriver() {
	ctx := context.Background()

	config, err := pgx.ParseConfig(configValues)
	if err != nil {
		log.Fatalf("parse config: %v\n", err)
	}
	config.Logger = &logger{}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("connect: %v\n", err)
	}

	sql := `SELECT schemaname, tablename FROM pg_tables WHERE schemaname = $1;`
	args := `information_schema`

	rows, err := conn.Query(ctx, sql, args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var pgtables []PgTable
	for rows.Next() {
		var s, t string
		if err := rows.Scan(&s, &t); err != nil {
			log.Fatal(err)
		}
		pgtables = append(pgtables, PgTable{
			SchemaName: s,
			TableName:  t})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, t := range pgtables {
		fmt.Println(t)
	}
}

var _ sqlhooks.Hooks = (*hook)(nil)

// implement Hooks interface
type hook struct{}

func (h *hook) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	log.Printf("SQL:\n%v\nARGS:%v\n", query, args)
	return ctx, nil
}

func (h *hook) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}

func loggingWithExtendedDriver() {
	ctx := context.Background()
	driverName := "postgres-proxy"
	sql.Register(driverName, sqlhooks.Wrap(stdlib.GetDefaultDriver(), &hook{}))
	db, err := sqlx.Connect(driverName, configValues)
	if err != nil {
		log.Fatalf("connect: %v\n", err)
	}

	var pgtables []PgTable
	sql := `SELECT schemaname, tablename FROM pg_tables WHERE schemaname = $1;`
	args := `information_schema`

	if err := db.SelectContext(ctx, &pgtables, sql, args); err != nil {
		log.Fatalf("select: %v\n", err)
	}

	for _, t := range pgtables {
		fmt.Println(t)
	}
}
