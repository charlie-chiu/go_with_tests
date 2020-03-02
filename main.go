package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	// open an empty file will get an error
	// can solve it by manuel add [] to the file
	if err != nil {
		log.Fatalf("could not opening %s %v", dbFileName, err)
	}
	store := NewFileSystemPlayerStore(db)
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
