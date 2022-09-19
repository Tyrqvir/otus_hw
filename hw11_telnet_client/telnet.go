package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	network    string
	connection net.Conn
}

func (client *client) Connect() error {
	connection, err := net.DialTimeout(client.network, client.address, client.timeout)
	if err != nil {
		return fmt.Errorf("error on connect: %w", err)
	}

	log.Println("...Connected to " + client.address)

	client.connection = connection

	return nil
}

func (client *client) Close() error {
	err := client.connection.Close()
	if err != nil {
		return fmt.Errorf("error on connection close: %w", err)
	}

	log.Println("...Connection was closed by peer")
	return nil
}

func (client *client) Send() error {
	_, err := io.Copy(client.connection, client.in)
	if err != nil {
		return fmt.Errorf("error on send: %w", err)
	}

	log.Println("...EOF")

	return nil
}

func (client *client) Receive() error {
	_, err := io.Copy(client.out, client.connection)
	if err != nil {
		return fmt.Errorf("error on receive: %w", err)
	}

	log.Println("receiving...")

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		network: "tcp",
	}
}
