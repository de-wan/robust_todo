// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db_sqlc

import (
	"database/sql"
	"time"
)

type Todo struct {
	Uuid       string
	Value      string
	CreatedAt  time.Time
	DoneAt     sql.NullTime
	Deadline   sql.NullTime
	ArchivedAt sql.NullTime
}
