package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/typisttech/wpsecadv/internal/wordfence"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()

	err := run(ctx, os.Args, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stderr io.Writer) error {
	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := mustParseFlags(args, stderr)

	timeoutCtx, cancel := context.WithTimeout(sigCtx, cfg.timeout)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(stderr, &slog.HandlerOptions{Level: cfg.level}))
	logger.Debug("Parsed config from flags", "config", cfg)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	s, err := newStore(logger, cfg.parallel, cwd, cfg.out)
	if err != nil {
		logger.Error("Failed to initialize data store", "error", err)
		return err
	}
	logger.Debug("Initialized data store", "store", s)

	g, gCtx := errgroup.WithContext(timeoutCtx)
	g.SetLimit(max(cfg.parallel, 4))

	rc, errc := wordfence.Fetch(gCtx, logger, nil, cfg.url, cfg.token)
	if rc == nil {
		return <-errc
	}

	g.Go(func() error {
		select {
		case e, ok := <-errc:
			if !ok {
				return nil
			}

			return e
		case <-gCtx.Done():
			return gCtx.Err()
		}
	})

	g.Go(func() error {
		for {
			select {
			case r, ok := <-rc:
				if !ok {
					return nil
				}

				g.Go(func() error {
					return s.Insert(r)
				})
			case <-gCtx.Done():
				return gCtx.Err()
			}
		}
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			err = context.Cause(ctx)
		}

		logger.Error("Failed to insert records", "error", err)
		return err
	}

	if err := s.Close(timeoutCtx); err != nil {
		if errors.Is(err, context.Canceled) {
			err = context.Cause(timeoutCtx)
		}

		logger.Error("Failed to close data store", "error", err)
		return err
	}

	return nil
}
