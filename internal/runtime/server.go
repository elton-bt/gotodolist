package runtime

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(server *http.Server, shutdownTimeout time.Duration, onReady func(addr string)) error {
	listener, err := net.Listen("tcp", listenAddress(server))
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	go func() {
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		close(errCh)
	}()

	if onReady != nil {
		onReady(listener.Addr().String())
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err, ok := <-errCh:
		if !ok {
			return nil
		}
		return err
	case <-signalCtx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

func listenAddress(server *http.Server) string {
	if server.Addr == "" {
		return ":http"
	}

	return server.Addr
}
