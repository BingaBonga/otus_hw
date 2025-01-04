package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalf("Usage: go-telnet %s %s", "host", "port")
	}

	ctx, closeCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer client.Close()

	go func() {
		defer closeCtx()
		err := client.Receive()
		if err != nil {
			log.Printf("cannot start client receive: %v\n", err)
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			log.Printf("cannot start client send: %v\n", err)
		}
	}()

	<-ctx.Done()
}
