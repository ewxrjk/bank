package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configPutCmd)
	configCmd.AddCommand(configGetCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "get configuration item",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return errors.New("usage: bank config get KEY")
		}
		if err = setup(); err != nil {
			return
		}
		var value string
		if value, err = b.GetConfig(args[0]); err != nil {
			return
		}
		fmt.Printf("%s\n", value)
		return
	},
}

var configPutCmd = &cobra.Command{
	Use:     "put KEY VALUE",
	Aliases: []string{"set"},
	Short:   "put configuration item",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 2 {
			return errors.New("usage: bank config put KEY VALUE")
		}
		if err = setup(); err != nil {
			return
		}
		if err = b.PutConfig(args[0], args[1]); err != nil {
			return
		}
		return
	},
}
