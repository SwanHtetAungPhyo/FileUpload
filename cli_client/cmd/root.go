/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/99designs/keyring"
	"os"

	"github.com/spf13/cobra"
)

var (
	distFile string
	key      []byte
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli_client",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCommand)
	rootCmd.AddCommand(statusCommand)
	rootCmd.AddCommand(commitCommand)
	rootCmd.AddCommand(remoteUrlCommand)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(KeyGenCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCommand.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	statusCommand.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	commitCommand.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	remoteUrlCommand.Flags().StringVarP(&remoteUrl, "url", "u", "", "remote url")
	remoteUrlCommand.Flags().StringVarP(&token, "token", "t", "", "token")
	uploadCmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	downloadCmd.Flags().StringVarP(&distFile, "dist", "d", "", "file path")

	KeyGenCmd.Flags().StringVarP(&username, "username", "u", "", "username")
	KeyGenCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "key-gen",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	keys, err := ring.Get("cipher-key")
	if err != nil {
		fmt.Println(err.Error())
	}
	key = keys.Data
}
