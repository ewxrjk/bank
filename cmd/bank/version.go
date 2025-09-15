package main

import (
	"fmt"

	"github.com/ewxrjk/bank/pkg/bank"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err = fmt.Printf("%s\n", bank.Version); err != nil {
			return
		}
		return
	},
}
