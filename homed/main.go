package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	done := make(chan bool)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		done <- true
	}()

	d := NewDaemon(&http.Client{
		Timeout: time.Second * 10,
	}, done)
	d.Run()
}
