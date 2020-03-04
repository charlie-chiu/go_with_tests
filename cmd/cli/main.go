package main

import (
	"fmt"
	poker "github.com/charlie-chiu/go_with_test"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, closeFunc, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFunc()

	fmt.Println("Let's play poker")
	fmt.Println("Type \"{Name} wins\" to record a win")

	//game := poker.CLI{store, os.Stdin}
	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}
