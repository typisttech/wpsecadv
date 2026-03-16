//go:build e2econtainer

package main

import (
	"bytes"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestE2E(t *testing.T) {
	ctx := t.Context()

	var out io.Writer
	if testing.Verbose() {
		out = t.Output()
	}

	srv, err := testcontainers.Run(
		ctx,
		"",
		testcontainers.WithName("wpsecadv-e2e-server"),
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{
			Context:        filepath.Join("..", ".."),
			Repo:           "localhost/wpsecadv-e2e-server",
			Tag:            "devel",
			BuildLogWriter: out,
			KeepImage:      true,
		}),
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForHTTP("/up").WithPort("8080/tcp"),
		),
	)
	t.Cleanup(func() {
		err := testcontainers.TerminateContainer(srv, testcontainers.StopTimeout(10*time.Second))
		if err != nil {
			t.Fatalf("failed to terminate server container: %v", err)
		}
	})
	if err != nil {
		t.Fatalf("failed to start server container: %v", err)
	}

	outerPort, err := srv.MappedPort(ctx, "8080/tcp")
	if err != nil {
		t.Fatalf("failed to get server container port: %v", err)
	}
	// TODO: The containers are in the same network.
	// There should be a way for them to connect without hopping to the host.
	// Maybe using compose? Maybe using custom networks? Maybe using the container name as the hostname?
	endpoint := "http://" + net.JoinHostPort(testcontainers.HostInternal, outerPort.Port())

	tests := []struct {
		name            string
		composerVersion string
	}{
		{"composer_2.9", "2.9"},
		{"composer_2.8", "2.8"},
		{"composer_2.7", "2.7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var out io.Writer
			if testing.Verbose() {
				out = t.Output()
			}

			scripts, err := testcontainers.Run(
				ctx,
				"",
				testcontainers.WithName("wpsecadv-e2e-scripts-"+tt.name),
				testcontainers.WithDockerfile(testcontainers.FromDockerfile{
					Context:        filepath.Join("..", ".."),
					Dockerfile:     filepath.Join("e2e", "scripts", "Dockerfile"),
					Repo:           "localhost/wpsecadv-e2e-scripts",
					Tag:            tt.name,
					BuildLogWriter: out,
					BuildArgs: map[string]*string{
						"COMPOSER_VERSION": &tt.composerVersion,
					},
					KeepImage: true,
				}),
				testcontainers.WithEntrypoint("sh", "-c", "echo ready && sleep infinity"),
				testcontainers.WithWaitStrategy(
					wait.ForLog("ready"),
				),
				testcontainers.WithHostPortAccess(outerPort.Int()), // TODO!
			)
			t.Cleanup(func() {
				err = testcontainers.TerminateContainer(scripts, testcontainers.StopTimeout(10*time.Second))
				if err != nil {
					t.Errorf("failed to terminate scripts container: %v", err)
				}
			})
			if err != nil {
				t.Fatalf("failed to start scripts container: %v", err)
			}

			cmd := []string{"/app/scripts.test"}
			if testing.Verbose() {
				cmd = append(cmd, "-test.v")
			}
			exec(t, scripts, cmd, tcexec.WithEnv([]string{"WPSECADV_SERVER_URL=" + endpoint}))
		})
	}
}

func exec(t *testing.T, c testcontainers.Container, cmd []string, opts ...tcexec.ProcessOption) {
	t.Helper()

	var buf bytes.Buffer
	buf.WriteString("### EXEC ")
	buf.WriteString(strings.Join(cmd, " "))

	code, out, err := c.Exec(t.Context(), cmd, opts...)
	if err != nil {
		t.Fatalf("%s\nfailed to exec command: %v", buf.String(), err)
	}

	buf.WriteRune('\n')
	_, err = buf.ReadFrom(out)
	if err != nil {
		t.Fatalf("error reading command output: %v", err)
	}

	if code != 0 {
		t.Fatalf("%s\nexit code = %d, want 0", buf.String(), code)
	}

	t.Log(buf.String())
}
