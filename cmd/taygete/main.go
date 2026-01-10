// Copyright (c) 2026 Michael D Henderson. All rights reserved.

// Package main implements the command line interface for the Olympia engine.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	addFlags := func(cmd *cobra.Command) error {
		cmd.PersistentFlags().Bool("debug", false, "enable debug logging (same as --log-level=debug)")
		cmd.PersistentFlags().Bool("quiet", false, "only log errors (same as --log-level=error)")
		cmd.PersistentFlags().String("log-level", "info", "logging level (debug|info|warn|error))")
		return nil
	}
	cmdRoot := &cobra.Command{
		Use:           "taygete",
		Short:         "taygete command line interface",
		Long:          `Command line interface for the Olympia game engine.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Root().PersistentFlags()
			logLevel, err := flags.GetString("log-level")
			if err != nil {
				return err
			}
			debug, err := flags.GetBool("debug")
			if err != nil {
				return err
			}
			quiet, err := flags.GetBool("quiet")
			if err != nil {
				return err
			}
			if debug && quiet {
				return fmt.Errorf("--debug and --quiet are mutually exclusive")
			}
			var lvl slog.Level
			switch {
			case debug:
				lvl = slog.LevelDebug
			case quiet:
				lvl = slog.LevelError
			default:
				switch strings.ToLower(logLevel) {
				case "debug":
					lvl = slog.LevelDebug
				case "info":
					lvl = slog.LevelInfo
				case "warn", "warning":
					lvl = slog.LevelWarn
				case "error":
					lvl = slog.LevelError
				default:
					return fmt.Errorf("log-level: unknown value %q", logLevel)
				}
			}
			handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     lvl,
				AddSource: lvl == slog.LevelDebug,
			})
			logger = slog.New(handler)
			slog.SetDefault(logger) // optional, but convenient
			return nil
		},
	}
	cmdRoot.AddCommand(cmdDb())
	cmdRoot.AddCommand(cmdVersion())
	err := addFlags(cmdRoot)
	if err != nil {
		logger.Error("root: addFlags",
			"err", err,
		)
		os.Exit(1)
	}

	err = cmdRoot.Execute()
	if err != nil {
		logger.Error("command failed", "err", err)
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
)
