package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stderr)

	timeout := flag.Duration("timeout", 10*time.Second, "connect to server timeout")
	flag.Parse()

	address, err := parseAddress(flag.Args())
	if err != nil {
		log.Fatalf("failed to parse arguments: %v", err)
	}

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("Connected to %s", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		defer cancel()
		runSender(client)
	}()

	go func() {
		defer cancel()
		runReceiver(client)
	}()

	<-ctx.Done()
}

func parseAddress(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("should set at 2 positional args")
	}

	host, port := args[0], args[1]
	if _, err := strconv.Atoi(port); err != nil {
		return "", errors.New("port should be integer value")
	}

	return net.JoinHostPort(host, port), nil
}

func runSender(client TelnetClient) {
	if err := client.Send(); err != nil {
		log.Println(err)
	}
	log.Println("...EOF")
}

func runReceiver(client TelnetClient) {
	if err := client.Receive(); err != nil {
		log.Println(err)
	}
	log.Println("...Connection was closed by peer")
}
