package test

import (
	"embed"
	"os"
)

var dsn = os.Getenv("POSTGRES_URL")

//go:embed *.sql
var migrations embed.FS
