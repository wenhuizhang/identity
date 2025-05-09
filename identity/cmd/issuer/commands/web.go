// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/agntcy/identity/pkg/httpserver"
	"github.com/agntcy/identity/pkg/log"

	"github.com/spf13/cobra"
)

const readHeaderTimeout = 10 * time.Second
const waitForServerDelay = 2 * time.Second
const waitForServerShutdownDelay = 5 * time.Second
const retryConnectionTimes = 10
const waitBeforeRetryDelay = 500 * time.Millisecond

var WebCmd = &cobra.Command{
	Use:   "web [port]",
	Short: "Starts the Web UI and keeps CLI active until Ctrl+C is pressed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the port is a valid number
		port := args[0]
		if _, err := strconv.Atoi(port); err != nil {
			log.Error("Invalid port number: ", port)
			return
		}

		// Check if the port is within the valid range
		portNum, _ := strconv.Atoi(port)
		if portNum < 1 || portNum > 65535 {
			log.Error("Port number out of range: ", port)
			return
		}

		// Check if the port is already in use
		if _, err := http.Get("http://localhost:" + port); err == nil {
			log.Error("Port already in use: ", port)
			return
		}

		log.Info("Starting Web UI...")

		// Create HTTP server with static file handler
		mux := http.NewServeMux()
		mux.Handle("/", httpserver.FileServer("web"))

		// Create an HTTP server with context
		srv := &http.Server{
			Addr:              ":" + port,
			Handler:           mux,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		// Channel to listen for OS signals (Ctrl+C)
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Run server in Goroutine
		go func() {
			log.Info("Web server running on http://localhost:" + port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error("Server error: ", err)
			}
		}()

		// Wait for the server to be ready before opening the browser
		go func() {
			time.Sleep(waitForServerDelay) // Short delay before checking
			waitForServer("http://localhost:" + port)
			openBrowser("http://localhost:" + port)
		}()

		// Wait for Ctrl+C signal
		<-stop
		log.Info("Shutting down server...")

		// Gracefully shutdown the server
		ctx, cancel := context.WithTimeout(context.Background(), waitForServerShutdownDelay)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Server forced to shutdown: ", err)
		}

		log.Info("Server exited.")
	},
}

// Waits for the server to be available
func waitForServer(url string) {
	ctx := context.Background()

	for range retryConnectionTimes { // Try X times before giving up
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		if err != nil {
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}

		resp.Body.Close()

		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("Web server is ready.")
			return
		}

		time.Sleep(waitBeforeRetryDelay) // Wait before retrying
	}

	log.Warn("Could not confirm server readiness. Try opening manually:", url)
}

// Opens the web browser
func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // Linux
		cmd = "xdg-open"
		args = []string{url}
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		log.Error("Failed to open browser: ", err)
	} else {
		log.Info("Browser opened successfully.")
	}
}
