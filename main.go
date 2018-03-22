package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		//err   error
		addr  string
		serve *Server
	)
	addr = fmt.Sprintf("%s:%d", "127.0.0.1", 8080)
	serve = NewServer(addr, 100, 100, 1024)
	go serve.Start()

	go func() {
		http.ListenAndServe("127.0.0.1:6060", nil)
	}()

	WaitSignal()
}

func WaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Printf("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		default:
			return
		}
	}
}
