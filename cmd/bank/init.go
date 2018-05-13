package main

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize bank database",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = setup(); err != nil {
			return
		}
		if err = b.NewBank(); err != nil {
			return
		}
		return
	},
}
