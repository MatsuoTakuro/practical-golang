package ch10

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func controlTranx() {
	db := openMyDB()
	tx := NewMiddlewareTx(db)

	http.Handle("/comments", tx(Recovery(http.HandlerFunc(Comments))))
	http.ListenAndServe(":8888", nil)

}

func openMyDB() *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=testuser dbname=testdb password=pass sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func NewMiddlewareTx(db *sql.DB) func(http.Handler) http.Handler {
	return func(wrappedHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, _ := db.Begin()
			lrw := NewLoggingResponseWriter(w)
			r = r.WithContext(context.WithValue(r.Context(), "tx", tx))

			wrappedHandler.ServeHTTP(lrw, r)

			// for Rollback
			// lrw.statusCode = http.StatusInternalServerError

			statusCode := lrw.statusCode // http.StatusOK (=200) as default
			if 200 <= statusCode && statusCode < 400 {
				log.Println("transaction committed")
				err := tx.Commit()
				if err != nil {
					lrw.WriteHeader(http.StatusInternalServerError)
				}
				_, err = lrw.Write([]byte("OK (in Response Body)\n"))
				if err != nil {
					lrw.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				log.Print("transaction rolling back due to status code: ", statusCode)
				err := tx.Rollback()
				if err != nil {
					lrw.WriteHeader(http.StatusInternalServerError)
				}
				_, err = lrw.Write([]byte("NG (in Response Body)\n"))
				if err != nil {
					lrw.WriteHeader(http.StatusInternalServerError)
				}
			}
		})
	}
}

func extractTx(r *http.Request) *sql.Tx {
	tx, ok := r.Context().Value("tx").(*sql.Tx)
	if !ok {
		panic("transaction middleware is not supported")
	}
	return tx
}

func Comments(w http.ResponseWriter, r *http.Request) {
	tx := extractTx(r)
	// DBアクセス処理
	log.Println(tx)
}
