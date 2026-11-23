package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"uuid"
)

type input struct {
	filename string
	duration time.Duration
}

func (inp *input) validate() error {
	inp.filename = strings.TrimSpace(inp.filename)
	if inp.filename == "" {
		return errors.New("no output file specified")
	}

	return nil
}

func die(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	var inp input
	flag.StringVar(&inp.filename, "file", "", "Output filename")
	flag.DurationVar(&inp.duration, "duration", 30*time.Second, "Log generation duration")
	flag.Parse()

	if err := inp.validate(); err != nil {
		die("invalid file: %s", err)
	}

	f, err := os.OpenFile(inp.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		die("cannot create output file: %s", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	timeout := time.NewTicker(inp.duration)
	defer timeout.Stop()

loop:
	for {
		temperature := rand.IntN(100) + 1
		instanceId := uuid.NewV4()
		logger.Info("recorded temperature", "temperature", temperature, "instance_id", instanceId)

		select {
		case <-ctx.Done():
			cancel()
			break loop

		case <-timeout.C:
			break loop

		default:
			// continue generating
		}
	}
}
