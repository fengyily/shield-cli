package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Process represents a running plugin process.
type Process struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
}

// StartPlugin launches a plugin binary and sends the start request.
// Returns the web port the plugin is listening on.
func StartPlugin(info *PluginInfo, cfg PluginConfig) (*Process, *StartResponse, error) {
	binPath := filepath.Join(PluginsDir(), info.Binary)

	// Verify binary exists
	if _, err := os.Stat(binPath); err != nil {
		return nil, nil, fmt.Errorf("plugin binary not found: %s", binPath)
	}

	cmd := exec.Command(binPath)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start plugin: %w", err)
	}

	proc := &Process{Cmd: cmd, Stdin: stdin, Stdout: stdout}

	// Send start request
	req := StartRequest{Action: "start", Config: cfg}
	if err := json.NewEncoder(stdin).Encode(req); err != nil {
		proc.Kill()
		return nil, nil, fmt.Errorf("failed to send start request: %w", err)
	}

	// Wait for ready response with timeout
	respCh := make(chan *StartResponse, 1)
	errCh := make(chan error, 1)
	go func() {
		var resp StartResponse
		if err := json.NewDecoder(stdout).Decode(&resp); err != nil {
			errCh <- fmt.Errorf("failed to read plugin response: %w", err)
			return
		}
		respCh <- &resp
	}()

	select {
	case resp := <-respCh:
		if resp.Status != "ready" {
			proc.Kill()
			return nil, nil, fmt.Errorf("plugin error: %s", resp.Message)
		}
		slog.Info("Plugin started", "name", resp.Name, "version", resp.Version, "web_port", resp.WebPort)
		return proc, resp, nil
	case err := <-errCh:
		proc.Kill()
		return nil, nil, err
	case <-time.After(15 * time.Second):
		proc.Kill()
		return nil, nil, fmt.Errorf("plugin did not respond within 15 seconds")
	}
}

// Stop sends a stop request to the plugin and waits for it to exit.
func (p *Process) Stop() {
	// Send stop request (best effort)
	req := StartRequest{Action: "stop"}
	json.NewEncoder(p.Stdin).Encode(req)
	p.Stdin.Close()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		p.Cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		slog.Warn("Plugin did not exit gracefully, killing")
		p.Cmd.Process.Kill()
	}
}

// Kill forcefully terminates the plugin process.
func (p *Process) Kill() {
	p.Stdin.Close()
	if p.Cmd.Process != nil {
		p.Cmd.Process.Kill()
	}
}
