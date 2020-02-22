package main

import (
	"flag"
	"fmt"
	"github.com/gohail/mafiosi/action"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"sync"
)

func main() {
	host := flag.String("host", "0.0.0.0:8080", "server's host:port")
	flag.Parse()
	logger := initZapLogger()
	defer logger.Sync()

	zap.S().Info("Mafiosi server run...")

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

func initZapLogger() *zap.Logger {
	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.DisableStacktrace = true
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger, _ := zapCfg.Build()
	zap.ReplaceGlobals(zapLogger)
	return zapLogger
}
