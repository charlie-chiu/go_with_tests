package main

import (
	"fmt"
	poker "github.com/charlie-chiu/go_with_test"
	"io"
	"log"
	"os"
	"time"
)

const dbFileName = "game.db.json"

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

	fmt.Println("Let's play poker")
	fmt.Println("Type \"{Name} wins\" to record a win")

	// just like HandlerFunc in http package
	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(Alerter), store)

	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
