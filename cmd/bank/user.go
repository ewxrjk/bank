package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var userPassword string

func init() {
	userAddCmd.PersistentFlags().StringVarP(&userPassword, "password", "p", "", "password")
	userPwCmd.PersistentFlags().StringVarP(&userPassword, "password", "p", "", "password")
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User administration",
}

var userAddCmd = &cobra.Command{
	Use:   "add USER",
	Short: "Add a user login",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return errors.New("usage: bank user add USER")
		}
		if err = setup(); err != nil {
			return
		}
		var pw string
		if pw, err = getPassword(); err != nil {
			return
		}
		if err = b.NewUser(args[0], pw); err != nil {
			return
		}
		return
	},
}
var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List user logins",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = setup(); err != nil {
			return
		}
		var users []string
		if users, err = b.GetUsers(); err != nil {
			return
		}
		for _, u := range users {
			if _, err = fmt.Println(u); err != nil {
				return
			}
		}
		return
	},
}

var userPwCmd = &cobra.Command{
	Use:   "pw USER",
	Short: "Set a user password",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return errors.New("usage: bank user pw USER")
		}
		if err = setup(); err != nil {
			return
		}
		var pw string
		if pw, err = getPassword(); err != nil {
			return
		}
		if err = b.SetPassword(args[0], pw); err != nil {
			return
		}
		return
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userPwCmd)
}

func getPassword() (password string, err error) {
	if userPassword != "" {
		password = userPassword
		return
	}
	var pw []byte
	var pw1, pw2 string
	for {
		fmt.Fprintf(os.Stderr, "Enter password: ")
		if pw, err = terminal.ReadPassword(0); err != nil {
			return
		}
		fmt.Fprintf(os.Stderr, "\n")
		pw1 = string(pw)

		fmt.Fprintf(os.Stderr, "Confirm password: ")
		if pw, err = terminal.ReadPassword(0); err != nil {
			return
		}
		fmt.Fprintf(os.Stderr, "\n")
		pw2 = string(pw)
		if pw1 == pw2 {
			password = pw1
			return
		}
		log.Printf("Passwords do not match - try again")
	}

}
