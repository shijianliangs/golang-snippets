package gracefulshutdown

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// RunHTTPServer runs srv.ListenAndServe in a goroutine and blocks until:
//   - ctx is canceled (then it triggers a graceful shutdown), or
//   - the server returns a non-ErrServerClosed error.
//
// shutdownTimeout bounds how long we wait for in-flight requests to finish.
func RunHTTPServer(ctx context.Context, srv *http.Server, shutdownTimeout time.Duration) error {
	errCh := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		// Trigger graceful shutdown.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		_ = srv.Shutdown(shutdownCtx)
		// Prefer the original ctx error (usually context.Canceled) only if the server didn't error.
		return <-errCh
	case err := <-errCh:
		return err
	}
}
