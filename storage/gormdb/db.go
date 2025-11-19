package gormdb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	Db    *gorm.DB
	SqlDb *sql.DB
}

func NewDb() *Db {
	return &Db{}
}

type DsnConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (d *Db) Open(config *DsnConfig) error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC search_path=public",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
	)
	var db *gorm.DB
	var err error

	for range 30 {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	d.Db = db
	d.SqlDb = sqlDb
	return nil
}

func (d *Db) Migrate(migrationPath string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	return goose.Up(d.SqlDb, migrationPath)
}
