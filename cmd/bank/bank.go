package main

import (
	"github.com/ewxrjk/bank"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = cobra.Command{
	Use:           "bank",
	Short:         "Household money tracking application.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

var dbDriver, dbSource string

func main() {
	rootCmd.PersistentFlags().StringVarP(&dbDriver, "driver", "D", "sqlite3", "database driver")
	rootCmd.PersistentFlags().StringVarP(&dbSource, "source", "d", "bank.db", "data source")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(configCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("bank: %s", err)
	}
}

var b *bank.Bank

func setup() (err error) {
	if b, err = bank.NewBank(dbDriver, dbSource); err != nil {
		return
	}
	return
}
