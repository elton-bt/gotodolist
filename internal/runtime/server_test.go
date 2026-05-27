package runtime

import (
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunInvokesReadyCallbackAfterSuccessfulBind(t *testing.T) {
	server := &http.Server{
		Addr: "127.0.0.1:0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	readyCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		errCh <- Run(server, time.Second, func(addr string) {
			readyCh <- addr
		})
	}()

	var addr string
	select {
	case addr = <-readyCh:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ready callback")
	}

	duplicateListener, err := net.Listen("tcp", addr)
	if err == nil {
		duplicateListener.Close()
		t.Fatalf("expected %q to stay bound while server is running", addr)
	}

	if err := server.Close(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("server.Close() returned error: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Run to return")
	}
}

func TestRunReturnsBindErrorBeforeReadyCallback(t *testing.T) {
	occupiedListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer occupiedListener.Close()

	readyCalled := false
	server := &http.Server{
		Addr: occupiedListener.Addr().String(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	err = Run(server, time.Second, func(string) {
		readyCalled = true
	})
	if err == nil {
		t.Fatal("expected bind error, got nil")
	}

	var opErr *net.OpError
	if !errors.As(err, &opErr) {
		t.Fatalf("expected net.OpError, got %T", err)
	}
	if readyCalled {
		t.Fatal("ready callback should not be called when bind fails")
	}
}
