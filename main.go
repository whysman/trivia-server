package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gameInfo struct {
	gameID   string
	userData []userEntry
}

type userEntry struct {
	name string
}

type userJoin struct {
	Name   string `json:"name"`
	GameID string `json:"gameid"`
}

type message struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

var activeGames = make(map[string]gameInfo)
var upgrader = websocket.Upgrader{}
var client *mongo.Client

func main() {
	client = connectToMongo()
	//http.HandleFunc("/createGame", createGame)
	//http.HandleFunc("/listGames", listGames)
	//http.HandleFunc("/joinGame", joinGame)

	http.HandleFunc("/wscomm", wsComm)
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func connectToMongo() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func initializeMongo() {
	collection := client.Database("trivia").Collection("users")
}

func wsComm(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		m := message{}
		json.Unmarshal(msg, &m)
		processMessage(m)
		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func getKeys(m map[string]gameInfo) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func processMessage(msg message) []byte {
	var response []byte
	if msg.Command == "createGame" {
		response = createGame()
	} else if msg.Command == "listGames" {
		response = listGames()
	} else if msg.Command == "joinGame" {
		response = joinGame(msg.Data)
	}
	return response
}

func createGame() []byte {
	name := namesgenerator.GetRandomName(0)
	activeGames[name] = gameInfo{gameID: name}
	jsonName, _ := json.Marshal(name)
	return jsonName
}

func listGames() []byte {
	jsonKeys, _ := json.Marshal(getKeys(activeGames))
	return jsonKeys
}

func joinGame(data string) []byte {
	b := []byte(data)
	u := userJoin{}
	json.Unmarshal(b, &u)
	if _, ok := activeGames[string(u.GameID)]; ok {
		ud := activeGames[string(u.GameID)].userData
		ud = append(ud, userEntry{name: u.Name})
		activeGames[string(u.GameID)] = gameInfo{gameID: string(u.GameID), userData: ud}
		fmt.Println(activeGames[string(u.GameID)].userData)
		return []byte("Joined game: " + data)
	}
	return []byte("Failed to join game: " + data)
}

/*
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
	u := userJoin{}
	json.Unmarshal(b, &u)
	if _, ok := activeGames[string(u.GameID)]; ok {
		ud := activeGames[string(u.GameID)].userData
		ud = append(ud, userEntry{name: u.Name})
		activeGames[string(u.GameID)] = gameInfo{gameID: string(u.GameID), userData: ud}
		fmt.Println(activeGames[string(u.GameID)].userData)
		w.WriteHeader(http.StatusOK)
	}
}
*/
