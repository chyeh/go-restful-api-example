package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	appCfg := newApplicationConfig()
	appCfg.bind(newEnvironmentVariableConfig())
	appCfg.bind(newCommandLineConfig())

	var apiCfg apiServerConfig
	apiCfg.load(appCfg)
	server := newAPIServer(apiCfg)
	go server.run()
	waitForSignal(server.shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func waitForSignal(shutdownFunc func(), signals ...os.Signal) {
	quitSig := make(chan os.Signal)
	signal.Notify(quitSig, signals...)
	fmt.Println(<-quitSig)
	shutdownFunc()
}
