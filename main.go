package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LJlamar/go_final_project/pkg/db"
	"github.com/LJlamar/go_final_project/pkg/server"
)

func main() {
	dbPath := os.Getenv("TODO_DBFILE")

	infoLogger := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

	//Checking set password
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		errorLogger.Fatalln("Variable TODO_PASSWORD isn't found or have empty value")
	}
	//Checking database
	if dbPath == "" {
		dbPath = "scheduler.db"
		infoLogger.Printf("Variable TODO_DBFILE isn't found, using default path: %v\n", dbPath)
	} else {
		infoLogger.Printf("Using path, given in TODO_DBFILE: %v\n", dbPath)
	}

	err := db.Init(dbPath)
	if err != nil {
		errorLogger.Fatalf("Failed to initialize database: %v", err)
	}
	//Connecting to server
	s := server.NewServer(infoLogger)
	infoLogger.Printf("Starting server at port %s\n", s.Server.Addr)

	err = http.ListenAndServe(s.Server.Addr, s.Server.Handler)
	if err != nil {
		errorLogger.Fatalln("Server failed: ", err)
	}
}
