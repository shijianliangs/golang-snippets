package gracefulshutdown

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunHTTPServer_ShutdownOnContextCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	srv := &http.Server{Handler: http.NewServeMux()}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- RunHTTPServer(ctx, srv, 2*time.Second)
	}()

	// Start serving on our listener.
	go func() {
		_ = srv.Serve(ln)
	}()

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for shutdown")
	}
}
