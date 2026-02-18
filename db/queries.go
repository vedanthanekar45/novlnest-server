package db

import "github.com/vedanthanekar45/novlnest-server/db/internal/database"

type Queries = database.Queries

func NewQueries(conn database.DBTX) *Queries {
	return database.New(conn)
}
