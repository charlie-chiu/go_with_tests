package poker

import (
	"io"
	"time"
)

type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}

// constructor of TexasHoldem
func NewTexasHoldem(a BlindAlerter, s PlayerStore) *TexasHoldem {
	return &TexasHoldem{
		alerter: a,
		store:   s,
	}
}

func (t *TexasHoldem) Start(numberOfPlayers int, dest io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	//blindIncrement := 5 * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		t.alerter.ScheduleAlertAt(blindTime, blind, dest)
		blindTime = blindTime + blindIncrement
	}
}

func (t *TexasHoldem) Finish(winner string) {
	t.store.RecordWin(winner)
}
