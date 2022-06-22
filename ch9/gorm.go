package ch9

import (
	"context"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PgTables struct {
	SchemaName string `gorm:"column:schemaname"`
	TableName  string `gorm:"column:tablename"`
}

func withGorm() {
	logger := Newlogger()
	db, err := gorm.Open(postgres.Open(configValues), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	var pgtables []PgTables
	ctx := context.Background()
	if err := db.WithContext(ctx).Where("schemaname = ?", "information_scheman").Find(&pgtables).Error; err != nil {
		log.Fatal(err)
	}

}

func Newlogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r", log.LstdFlags),
		logger.Config{LogLevel: logger.Info},
	)
}
