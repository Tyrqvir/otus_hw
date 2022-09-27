package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
)

type Args struct {
	Host string
	Port string
}

type Options struct {
	Timeout time.Duration `short:"t" long:"timeout" description:"Timeout" default:"10s" required:"yes"`
	Args    Args          `positional-args:"yes" required:"yes"`
}

var options Options

var errHelp flags.ErrorType

func parser() {
	if _, err := flags.NewParser(&options, flags.Default).Parse(); err != nil {
		if errors.As(err, &errHelp) {
			if errors.Is(err, flags.ErrHelp) {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
}

func main() {
	parser()

	address := net.JoinHostPort(options.Args.Host, options.Args.Port)
	ctx, cancel := context.WithCancel(context.Background())

	client, err := NewTelnetClient(address, options.Timeout, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	go func() {
		defer cancel()

		log.Println("go send")
		err := client.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer cancel()

		log.Println("go receive")

		err := client.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	signChannel := make(chan os.Signal, 1)
	signal.Notify(signChannel, syscall.SIGINT)

	select {
	case <-signChannel:
		cancel()
	case <-ctx.Done():
		close(signChannel)
	}
}
