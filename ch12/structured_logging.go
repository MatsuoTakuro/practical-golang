package ch12

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func structured_logging() {
	// simple()
	// fundamental()
	// otherInfo()
	errorcode()
}

func simple() {
	log.Printf("Hello World")
}

func fundamental() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Error().Msg("error in output")
	log.Info().Msg("info in output")
	log.Debug().Msgf("debug in output: %v", "however, no output")
}

func otherInfo() {
	log.Info().
		Str("app", "awesome-app").
		Int("user_id", 1114).
		Send()

	ctx := context.Background()
	logger := log.With().
		Int("user_id", 1024).
		Str("path", "/api/user").
		Str("method", "post").
		Logger()

	ctx = logger.WithContext(ctx)
	newLogger := zerolog.Ctx(ctx)
	newLogger.Print("debug message")
	newLogger.Print(ctx)

	logger2 := log.With().
		Int("user_id", 1025).
		Str("path", "/api/user").
		Str("method", "post").
		Logger()

	ctx = logger2.WithContext(ctx)
	newLogger2 := zerolog.Ctx(ctx)
	newLogger2.Print("debug message")
	newLogger2.Print(ctx)

}

const DBConnErrCD = "E10001"

func errorcode() {
	_, err := openDB()
	if err != nil {
		log.Error().
			Str("code", DBConnErrCD).
			Err(err).
			Msg("db connection failed")
	}
}

func openDB() (*sql.DB, error) {
	return nil, errors.New("on purpose")
}
