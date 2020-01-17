package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

var u = websocket.Upgrader{}

func main() {
	fmt.Println("Mafiosi server run...")

	logger, _ := zap.NewDevelopment()

	zap.ReplaceGlobals(logger)

	http.HandleFunc("/", stream)

	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startService(":8080", wg)
	zap.S().Info("Started on localhost:8080")
	zap.S().Info("WS handler: ws://localhost:8080/")
	wg.Wait()
}

func must(err error, description string) {
	if err == nil {
		return
	}

	if l := zap.S(); l != nil {
		// Use logger if available
		l.Fatalf("%s - %s", description, err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "%s - %s", description, err)
	os.Exit(1)
}

func startService(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	must(http.ListenAndServe(port, nil), "failed to start server")
}



