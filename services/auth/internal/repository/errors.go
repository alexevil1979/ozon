package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrDuplicateEmail = errors.New("email already exists")

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
