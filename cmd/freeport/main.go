package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	port, err := freePort()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, port)
}

// freePort asks the kernel for an unused TCP port and returns it. It backs the
// mise tasks that need a free port, replacing a shell-specific probe so they
// behave the same under bash and zsh on macOS and Ubuntu.
func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("reserving a free port: %w", err)
	}
	defer func() { l.Close() }()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("unexpected listener address type %T", l.Addr())
	}

	return addr.Port, nil
}
