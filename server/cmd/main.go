package main

import (
	"fmt"
	"github.com/psykhi/pong/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := server.NewServer()
	s.Start()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		s.Stop()
		os.Exit(0)
	}()
}
