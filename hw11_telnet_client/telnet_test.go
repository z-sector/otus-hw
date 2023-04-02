package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in, out := &bytes.Buffer{}, &bytes.Buffer{}

			client := NewTelnetClient(l.Addr().String(), 2*time.Second, io.NopCloser(in), out)
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

	t.Run("timeout case", func(t *testing.T) {
		in, out := &bytes.Buffer{}, &bytes.Buffer{}

		timeout := 2 * time.Second
		client := NewTelnetClient("google.com:36063", timeout, io.NopCloser(in), out)

		start := time.Now()
		err := client.Connect()
		end := time.Now()

		assert.InDelta(t, timeout.Seconds(), end.Sub(start).Seconds(), 0.5, fmt.Sprintf("Should timeout after %v", timeout))
		require.Error(t, err, "Timeout error expected")
		defer func() {
			require.NoError(t, client.Close())
		}()
	})

	t.Run("connection refused", func(t *testing.T) {
		in, out := &bytes.Buffer{}, &bytes.Buffer{}

		client := NewTelnetClient("localhost:36063", 2*time.Second, io.NopCloser(in), out)
		defer func() {
			require.NoError(t, client.Close())
		}()

		err := client.Connect()
		require.Error(t, err)
	})

	t.Run("send without connect", func(t *testing.T) {
		in, out := &bytes.Buffer{}, &bytes.Buffer{}

		client := NewTelnetClient("localhost:36063", 2*time.Second, io.NopCloser(in), out)
		defer func() {
			require.NoError(t, client.Close())
		}()

		err := client.Send()
		require.Equal(t, ErrNeedToConnect, err)
	})

	t.Run("receive without connect", func(t *testing.T) {
		in, out := &bytes.Buffer{}, &bytes.Buffer{}

		client := NewTelnetClient("localhost:36063", 2*time.Second, io.NopCloser(in), out)
		defer func() {
			require.NoError(t, client.Close())
		}()

		err := client.Receive()
		require.Equal(t, ErrNeedToConnect, err)
	})

	t.Run("close receiver", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, l.Close())
		}()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in, out := &bytes.Buffer{}, &bytes.Buffer{}

			client := NewTelnetClient(l.Addr().String(), 2*time.Second, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			require.Eventually(t, func() bool {
				in.WriteString("hello\n")
				err := client.Send()
				return err != nil
			}, time.Second, time.Millisecond)
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() {
				require.NoError(t, conn.Close())
			}()

			n, err := conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}
