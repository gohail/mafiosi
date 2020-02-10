package main

import (
	"flag"
	"fmt"
	"github.com/gohail/mafiosi/action"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

func main() {
	host := flag.String("host", "localhost:8080", "server's host:port")
	flag.Parse()

	fmt.Println("Mafiosi server run...")

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	http.HandleFunc("/", action.ConnHandler)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startService(*host, wg)
	zap.S().Info("Started on ", *host)
	zap.S().Infof("WS handler: ws://%s/", *host)
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

func startService(host string, wg *sync.WaitGroup) {
	defer wg.Done()
	must(http.ListenAndServe(host, nil), "failed to start server")
}
