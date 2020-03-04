package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in    *bufio.Scanner
}

//just like constructor in php?
func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		store: store,
		in:    bufio.NewScanner(in),
	}
}

func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.store.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.TrimSuffix(userInput, " wins")
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
