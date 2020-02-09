package main

import (
	"errors"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configPutCmd)
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

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "list current configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = setup(); err != nil {
			return
		}
		var config map[string]string
		if config, err = b.GetConfigs(); err != nil {
			return
		}
		keys := make([]string, 0, len(config))
		for k := range config {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if _, err = fmt.Printf("%s=%s\n", k, config[k]); err != nil {
				return
			}
		}
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
