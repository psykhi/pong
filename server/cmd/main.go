package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/psykhi/pong/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var conf server.Config
	envconfig.MustProcess("", &conf)
	s := server.NewServer(conf)
	s.Start()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Terminating sever...")
		s.Stop()
		os.Exit(0)
	}()
}
