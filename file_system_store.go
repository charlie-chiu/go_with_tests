package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func NewFileSystemPlayerStore(db *os.File) *FileSystemPlayerStore {
	db.Seek(0, 0)
	league, err := NewLeague(db)
	if err != nil {
		log.Fatalf("new league failed, %v", err)
	}
	return &FileSystemPlayerStore{
		database: &tape{db},
		league:   league,
	}
}

type FileSystemPlayerStore struct {
	database io.Writer
	league   League
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	p := f.league.find(name)
	if p != nil {
		return p.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	json.NewEncoder(f.database).Encode(f.league)
}

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}
