package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"

	// source/file import is required for migration files to read
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	// load pq as database driver
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Ved1234"
	dbname   = "gRPC"
)

var DB *sqlx.DB

func DbConnection() (*sqlx.DB, error) {
	postgresql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", postgresql)
	if err != nil {
		fmt.Println("err")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("database connection is done")

	driver, driverErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driverErr != nil {
		return nil, driverErr
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres://postgres:Ved1234@localhost:5432/gRPC?sslmode=disable&search_path=public",
		driver)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error in migration ", err)
		return nil, err
	}

	return db, nil
}
