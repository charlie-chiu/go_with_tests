package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (FileSystemPlayerStore) GetPlayerScore(name string) int {
	panic("implement me")
}

func (FileSystemPlayerStore) RecordWin(name string) {
	panic("implement me")
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
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
