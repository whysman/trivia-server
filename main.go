package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	http.HandleFunc("/listGames", listGames)
	http.HandleFunc("/joinGame", joinGame)
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
	jsonName, _ := json.Marshal(name)
	w.Write([]byte(jsonName))
	return
}

func listGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	jsonKeys, _ := json.Marshal(getKeys(activeGames))
	w.Write([]byte(jsonKeys))
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
}

func getKeys(m map[string]gameInfo) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
