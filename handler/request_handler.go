package handler

import "github.com/jmoiron/sqlx"

type RequestHandler struct {
	db *sqlx.DB
}
