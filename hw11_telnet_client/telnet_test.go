package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var address = "google.ru:80"

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client, err := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, err)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("correct instance of client", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		address := "google.ru:80"
		timeout := time.Duration(10)

		expectedClient := &client{
			address: address,
			timeout: timeout,
			in:      ioutil.NopCloser(in),
			out:     out,
			network: "tcp",
		}
		client, err := NewTelnetClient(address, timeout, ioutil.NopCloser(in), out)
		require.NoError(t, err)

		require.Equal(t, expectedClient, client)
	})

	t.Run("failed connection - not found port", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		address := "google.ru"
		timeout := time.Duration(10)

		client, err := NewTelnetClient(address, timeout, ioutil.NopCloser(in), out)
		require.NoError(t, err)

		err = client.Connect()
		require.Equal(t, "error on connect: dial tcp: address google.ru: missing port in address", err.Error())
	})

	t.Run("failed client - not found reader", func(t *testing.T) {
		out := &bytes.Buffer{}
		timeout := time.Duration(10)

		_, err := NewTelnetClient(address, timeout, nil, out)
		require.Equal(t, "reader not defined", err.Error())
	})

	t.Run("failed client - not found writer", func(t *testing.T) {
		in := &bytes.Buffer{}
		timeout := time.Duration(10)

		_, err := NewTelnetClient(address, timeout, ioutil.NopCloser(in), nil)
		require.Equal(t, "writer not defined", err.Error())
	})
}
