package ch9

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func rollbackWithDefer() {
	s := Service{
		db: db,
	}
	ctx := context.Background()
	err := s.UpdateProduct(ctx, "0001")
	if err != nil {
		log.Fatalf("update product; %v", err)
	}
}

func seperateTxCtrlAndImpl() {
	s := Service2{
		tx: txAdmin{
			db: db,
		},
	}
	ctx := context.Background()
	err := s.UpdateProductOnly(ctx, "0001")
	if err != nil {
		log.Fatalf("update product; %v", err)
	}
}

type Service struct {
	db *sql.DB
}

func (s *Service) UpdateProduct(ctx context.Context, productID string) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	cmd := `UPDATE products SET price = 200 WHERE product_id = $1;`
	if _, err = tx.ExecContext(ctx, cmd, productID); err != nil {
		return err
	}

	return tx.Commit()
}

type txAdmin struct {
	db *sql.DB
}

type Service2 struct {
	tx txAdmin
}

func (t *txAdmin) Transaction(ctx context.Context, f func(ctx context.Context) (err error)) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := f(ctx); err != nil {
		return fmt.Errorf("transaction query failed: %w", err)
	}
	return tx.Commit()
}

func (s *Service2) UpdateProductOnly(ctx context.Context, productID string) (err error) {
	updateFunc := func(ctx context.Context) error {
		cmd := `UPDATE products SET price = 300 WHERE product_id = $1;`
		if _, err = s.tx.db.ExecContext(ctx, cmd, productID); err != nil {
			return err
		}
		return nil
	}
	return s.tx.Transaction(ctx, updateFunc)
}
