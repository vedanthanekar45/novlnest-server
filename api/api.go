package api

import (
	"github.com/vedanthanekar45/novlnest-server/db"
	"golang.org/x/oauth2"
)

type ApiConfig struct {
	Store *db.Queries
	OAuth *oauth2.Config
}
