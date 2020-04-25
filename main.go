package main

import (
	"fmt"
	"net/http"

	"github.com/docker/docker/pkg/namesgenerator"
)

type gameInfo struct {
	gameID   string
	userData []userEntry
}

type userEntry struct {
	name string
}

var activeGames = make(map[string]gameInfo)

func main() {
	http.HandleFunc("/createGame", createGame)
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println(err)
	}

}

func createGame(w http.ResponseWriter, r *http.Request) {
	name := namesgenerator.GetRandomName(0)
	activeGames[name] = gameInfo{gameID: name}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(name))
	return
}
