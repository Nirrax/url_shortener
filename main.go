package main

import "log"

func main() {
	db, err := NewPostgresDB()

	if err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":3000", db)
	server.Run()
}
