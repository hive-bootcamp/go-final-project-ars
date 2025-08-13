package main

import (
	"log"

	"github.com/Arsadidas/go-final-project/pkg/db"
	"github.com/Arsadidas/go-final-project/pkg/server"
	"github.com/Arsadidas/go-final-project/tests"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := db.Init(tests.DBFile)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return
	}
	server.Start()
}
