package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/rom8726/etoggl/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "app",
}

func main() {
	rootCmd.AddCommand(server.ServerCmd)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1) //nolint:gocritic // it's ok here
	}
}
