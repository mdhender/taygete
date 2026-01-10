// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"log"

	"github.com/mdhender/taygete"
	"github.com/spf13/cobra"
)

func cmdVersion() *cobra.Command {
	showBuildInfo := false
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().BoolVar(&showBuildInfo, "build-info", showBuildInfo, "show build information")
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "display the application's version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showBuildInfo {
				fmt.Println(taygete.Version().String())
				return nil
			}
			fmt.Println(taygete.Version().Core())
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}
