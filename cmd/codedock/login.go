package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"codedock/cmd/codedock/client"
	"codedock/cmd/codedock/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your Codedock server",
	Long:  `Authenticate your CLI with a self-hosted Codedock server instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Codedock Server URL (e.g. https://api.yourdomain.com): ")
		serverURL, _ := reader.ReadString('\n')
		serverURL = strings.TrimSpace(serverURL)
		serverURL = strings.TrimSuffix(serverURL, "/")

		if serverURL == "" {
			fmt.Println("Error: Server URL is required.")
			os.Exit(1)
		}

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println("Error reading password:", err)
			os.Exit(1)
		}
		password := string(bytePassword)

		fmt.Println("Authenticating...")

		client := client.NewClient(serverURL, "")

		authResp, err := client.Login(email, password)
		if err != nil && strings.Contains(err.Error(), "2FA code required") {
			fmt.Print("Enter 2FA Code: ")
			byteTotp, readErr := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if readErr != nil {
				fmt.Printf("❌ Failed to read 2FA code: %v\n", readErr)
				os.Exit(1)
			}
			authResp, err = client.LoginWithTOTP(email, password, strings.TrimSpace(string(byteTotp)))
		}
		if err != nil {
			fmt.Printf("❌ Authentication failed: %v\n", err)
			os.Exit(1)
		}

		if err := config.Save(&config.Config{
			ServerURL: serverURL,
			Token:     authResp.Token,
			Email:     authResp.User.Email,
		}); err != nil {
			fmt.Printf("❌ Failed to save configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully logged in as %s\n", authResp.User.Email)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
