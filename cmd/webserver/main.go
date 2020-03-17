package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	poker "github.com/charlie-chiu/go_with_test"
)

const dbFileName = "game.db.json"
const port = 5000

func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}

func main() {
	store, closeFunc, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFunc()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(Alerter), store)
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatalf("could not create new server %v", err)
	}

	fmt.Fprintf(os.Stdout, "start listen on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
