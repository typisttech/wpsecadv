package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/typisttech/wpsecadv/internal/wordfence/vuln"
)

type config struct {
	level slog.Level

	parallel int
	timeout  time.Duration

	out string

	url   string
	token vuln.Token
}

func mustParseFlags(args []string, stderr io.Writer) config {
	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	l := flags.String(
		"loglevel",
		"INFO",
		"Log `level`. Valid values are DEBUG, INFO, WARN, ERROR.",
	)

	flags.IntVar(
		&cfg.parallel,
		"parallel",
		0,
		"run `n` workers simultaneously."+`
If n is 0 or less, GOMAXPROCS is used. Setting -parallel to values higher
 than GOMAXPROCS may cause degraded performance due to CPU contention.
(default GOMAXPROCS)`)

	flags.DurationVar(
		&cfg.timeout,
		"timeout",
		5*time.Minute,
		"If data generation runs longer than duration `d`, abort.",
	)

	// For data store.
	flags.StringVar(
		&cfg.out,
		"out",
		"assets",
		"Path to assets `directory`."+`
If the directory already exists, it will be emptied.`,
	)

	flags.StringVar(
		&cfg.url,
		"url",
		vuln.FeedProduction,
		"URL of the Wordfence vulnerability data `feed`.",
	)

	flags.Var(
		&cfg.token,
		"token",
		"Wordfence intelligence API Key to send as Authorization: Bearer <`token`>.",
	)

	// Ignore error because of flag.ExitOnError
	_ = flags.Parse(args[1:])

	var level slog.Level
	if err := level.UnmarshalText([]byte(*l)); err != nil {
		fmt.Fprintf(stderr, "invalid value %q for flag -loglevel: %v\n", *l, err)
		flags.Usage()
		os.Exit(2)
	}
	cfg.level = level

	if cfg.parallel < 1 {
		cfg.parallel = runtime.GOMAXPROCS(0)
	}

	if cfg.out == "" {
		fmt.Fprintln(stderr, "invalid value for flag -out: cannot be empty")
		flags.Usage()
		os.Exit(2)
	}

	if cfg.url == "" {
		fmt.Fprintln(stderr, "default -url to production feed")
		cfg.url = vuln.FeedProduction
	}

	return cfg
}

func (c config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("level", c.level.String()),
		slog.Int("parallel", c.parallel),
		slog.Duration("timeout", c.timeout),
		slog.String("out", c.out),
		slog.String("url", c.url),
		slog.Any("token", c.token),
	)
}
