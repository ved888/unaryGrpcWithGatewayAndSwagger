package dbhelper

import "github.com/jmoiron/sqlx"

type DAO struct {
	DB *sqlx.DB
}
