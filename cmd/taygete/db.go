// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mdhender/taygete"
	"github.com/spf13/cobra"
)

func cmdDb() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "db",
		Short: "database commands",
	}
	cmd.AddCommand(cmdDbInit())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdDbInit() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "initialize a new database",
		Args:  cobra.MinimumNArgs(1), // path to database
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			if isdir(path) {
				path = filepath.Join(path, "taygete.db")
			}
			dbPath, dbName := filepath.Dir(path), filepath.Base(path)
			if !isdir(dbPath) {
				err := fmt.Errorf("path does not exist: %q", dbPath)
				logger.Error("db: init",
					"err", err)
				return err
			}
			if filepath.Ext(dbName) != ".db" {
				err := fmt.Errorf("name must have '.db' suffix: %q", dbName)
				logger.Error("db: init",
					"err", err)
				return err
			}
			if isfile(path) {
				err := fmt.Errorf("database exists: %q", path)
				logger.Error("db: init",
					"err", err)
				return err
			}
			db, err := taygete.OpenGameDB(path)
			if err != nil {
				logger.Error("db: init",
					"err", err)
			}
			defer func() { _ = db.Close() }()
			logger.Info("db: init",
				"created", path)
			teg, err := taygete.NewEngine(db, nil)
			if err != nil {
				logger.Error("db: init",
					"err", err)
			}
			_ = teg
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func isdir(path string) bool {
	sb, err := os.Stat(path)
	if err != nil {
		return false
	}
	return sb.IsDir()
}

func isfile(path string) bool {
	sb, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !sb.IsDir() && sb.Mode().IsRegular()
}
