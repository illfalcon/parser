package main

import (
	"github.com/illfalcon/parser/internal/db"
	"github.com/illfalcon/parser/internal/frontend"
)

func main() {
	db.Prepare()
	frontend.Start()
}
