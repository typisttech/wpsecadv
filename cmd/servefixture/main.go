package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run serves the given directory over HTTP. It backs the gen:data:fixture
// mise task, replacing an external static file server with the standard
// library so development needs one fewer dependency.
func run(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("usage: %s <addr> <directory>", args[0])
	}

	srv := &http.Server{
		Addr:              args[1],
		Handler:           http.FileServer(http.Dir(args[2])),
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("serving %q: %w", args[2], err)
	}

	return nil
}
