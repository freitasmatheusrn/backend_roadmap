/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/freitasmatheusrn/caching-proxy/server"
	"github.com/spf13/cobra"
)

// cachingProxyCmd represents the cachingProxy command
var cachingProxyCmd = &cobra.Command{
	Use:   "caching-proxy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
			and usage of using your command. For example:

			Cobra is a CLI library for Go that empowers applications.
			This application is a tool to generate the needed files
			to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := url.ParseRequestURI(origin); err != nil {
			return fmt.Errorf("flag --origin needs to be valid URL: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		server := server.NewServer(port, origin)
		fmt.Printf("running proxy on http://localhost:%s -> %s\n", port, origin)
		go func() {
			if err := server.InitServer(); err != nil {
				fmt.Printf("Error initialing server: %s\n", err)
			}
		}()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("\n> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("error reading command: %s\n", err)
				continue
			}
			command := strings.TrimSpace(input)

			switch command {
			case "clear-cache":
				server.ClearCache()
				fmt.Println("🧹 Cache cleared.")
			case "status":
				fmt.Println("🔍 Server running...")
			case "exit", "quit":
				fmt.Println("👋 shutting down server...")
				os.Exit(0)
			default:
				fmt.Println("❓ Command not recognized. Available commands: clear-cache, status, exit")
			}
		}
	},
}
var port string
var origin string

func init() {
	rootCmd.AddCommand(cachingProxyCmd)
	cachingProxyCmd.Flags().StringVarP(&port, "port", "p", "", "caching proxy server port")
	cachingProxyCmd.Flags().StringVarP(&origin, "origin", "o", "", "server to forward the requests to")
	cachingProxyCmd.Flags().Bool("clear-cache", false, "clear the server cache")
	cachingProxyCmd.MarkFlagsRequiredTogether("port", "origin")
}
