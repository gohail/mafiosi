package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
)

var wsCount int

// Accessory app for join players in game
func main() {
	clCount := flag.Int("n", 1, "enter number of game's clients which will be created")
	gameId := flag.Int("id", 0, "enter game id for connection")
	flag.Parse()

	if *gameId == 0 {
		fmt.Println("game id is not defined")
		os.Exit(1)
	}

	fmt.Println("create ws clients n:", *clCount)
	var wsSli []*websocket.Conn
	for i := 0; i < *clCount; i++ {
		wsSli = append(wsSli, createWSClient())
	}

	// join all
	joinAllCli(wsSli)
	printAllReq(wsSli)

}

func createWSClient() *websocket.Conn {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/",
	}
	wsCount++
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	fmt.Printf(" created WS# %d client url: %s \n", wsCount, u.String())
	return c
}

func joinAllCli(wsSli []*websocket.Conn) {
	sendJoin := []byte(`{"action" : "JOIN_GAME"}`)
	var err error
	for i, v := range wsSli {
		fmt.Printf("WS:%d send JOIN \n", i)
		err = v.WriteMessage(1, sendJoin)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func printAllReq(wsSli []*websocket.Conn) {
	for i, v := range wsSli {
		_, p, err := v.ReadMessage()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("WS:%d got: %s", i, string(p))
	}
}

func typeGameId(id int, wsSli []*websocket.Conn) {

}
