package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/typisttech/wpsecadv/internal/data"
	"github.com/typisttech/wpsecadv/internal/server"
)

const (
	defaultPort          = "8080"
	defaultShutdownLimit = 8 * time.Second
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		fmt.Fprintln(os.Stdout, "Exiting")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "Exiting")
}

func run(ctx context.Context, w io.Writer) error {
	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := log.New(w, "", log.LstdFlags|log.LUTC)
	printInfo(logger)

	port, err := port()
	if err != nil {
		return err
	}

	shutdownTimeout, err := gracefulShutdownTimeout()
	if err != nil {
		return err
	}

	modTime := vcsTimeOrNow()

	logger.Printf(" * %-15s%s", "304 Mod Time:", modTime.Format(time.RFC3339))

	srv := &http.Server{
		Addr:        ":" + port,
		Handler:     server.New(&data.Store{}, modTime),
		ReadTimeout: 5 * time.Second, // TODO: Allow customization.
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Printf(" * Listening on http://127.0.0.1%s", srv.Addr)
		logger.Printf(" * Listening on http://[::1]%s", srv.Addr)
		logger.Println("Use Ctrl-C to stop")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-sigCtx.Done():
		stop()
		logger.Printf("Stopping the server: %v", context.Cause(sigCtx))
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	logger.Println("Gracefully stopping, waiting for requests to finish")
	if err := srv.Shutdown(shutdownCtx); err != nil { //nolint: contextcheck
		return err
	}

	logger.Println("Server has stopped")
	return nil
}

func port() (string, error) {
	p := os.Getenv("PORT")
	if p == "" {
		return defaultPort, nil
	}

	i, err := strconv.Atoi(p)
	if err != nil {
		return "", fmt.Errorf("invalid PORT: %q is not an integer", p)
	}
	if i <= 0 {
		return "", fmt.Errorf("invalid PORT: %d is not a positive integer", i)
	}

	return p, nil
}

func gracefulShutdownTimeout() (time.Duration, error) {
	n := os.Getenv("WP_SEC_ADV_GRACEFUL_SHUTDOWN_TIMEOUT")
	if n == "" {
		return defaultShutdownLimit, nil
	}

	i, err := strconv.Atoi(n)
	if err != nil {
		return 0, fmt.Errorf("invalid WP_SEC_ADV_GRACEFUL_SHUTDOWN_TIMEOUT: %q is not an integer", n)
	}
	if i <= 0 {
		return 0, fmt.Errorf("invalid WP_SEC_ADV_GRACEFUL_SHUTDOWN_TIMEOUT: %d is not a positive integer", i)
	}

	return time.Duration(i) * time.Second, nil
}

func vcsTimeOrNow() time.Time {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return time.Now().UTC()
	}

	for _, s := range bi.Settings {
		if s.Key != "vcs.time" {
			continue
		}

		t, err := time.Parse(time.RFC3339, s.Value)
		if err != nil {
			break
		}

		return t.UTC()
	}

	return time.Now().UTC()
}

func printInfo(logger *log.Logger) {
	logger.Println("==> Booting WP Sec Adv")

	keys := map[string]string{
		"GOARCH":       "Go Arch",
		"GOOS":         "Go OS",
		"vcs.revision": "VCS Revision",
		"vcs.time":     "VCS Time",
		"vcs.modified": "VCS Dirty",
	}
	bi, ok := debug.ReadBuildInfo()
	if ok {
		logger.Printf(" * %-15s%s", "Go Version:", bi.GoVersion)

		for _, s := range bi.Settings {
			label, ok := keys[s.Key]
			if !ok {
				continue
			}

			logger.Printf(" * %-15s%s", label+":", s.Value)
		}
	}
}
