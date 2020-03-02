package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	p := f.GetLeague().find(name)
	if p != nil {
		return p.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}

func (f *FileSystemPlayerStore) GetLeague() League {
	f.database.Seek(0, 0)
	league, err := NewLeague(f.database)
	if err != nil {
		log.Fatal(err)
	}

	return league
}

func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}
