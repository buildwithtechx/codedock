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

var (
	loginServerURL string
	loginToken     string
	loginTokenFile string
	loginEmail     string
	loginPassword  string
	loginLocal     bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your Codedock server",
	Long:  `Authenticate your CLI with a self-hosted Codedock server instance using email/password, Personal Access Token, or local operator credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverURL := strings.TrimSpace(loginServerURL)
		if serverURL == "" {
			serverURL = os.Getenv("CODEDOCK_SERVER_URL")
		}
		if serverURL == "" {
			if existing, err := config.Load(); err == nil && existing != nil && existing.ServerURL != "" {
				serverURL = existing.ServerURL
			}
		}

		reader := bufio.NewReader(os.Stdin)

		if serverURL == "" {
			fmt.Print("Codedock Server URL (e.g. https://api.yourdomain.com): ")
			input, _ := reader.ReadString('\n')
			serverURL = strings.TrimSpace(input)
		}
		serverURL = strings.TrimSuffix(serverURL, "/")

		if serverURL == "" {
			fmt.Println("Error: Server URL is required.")
			os.Exit(1)
		}

		token := strings.TrimSpace(loginToken)
		if token == "" {
			token = strings.TrimSpace(os.Getenv("CODEDOCK_TOKEN"))
		}
		if token == "" && loginTokenFile != "" {
			data, err := os.ReadFile(loginTokenFile)
			if err != nil {
				fmt.Printf("❌ Failed to read token file: %v\n", err)
				os.Exit(1)
			}
			token = strings.TrimSpace(string(data))
		}

		apiClient := client.NewClient(serverURL, "")

		if loginLocal {
			internalToken := os.Getenv("CODEDOCK_INTERNAL_TOKEN")
			if internalToken == "" {
				fmt.Print("Internal Operator Token: ")
				byteToken, err := term.ReadPassword(int(syscall.Stdin))
				fmt.Println()
				if err != nil {
					fmt.Printf("❌ Failed to read internal token: %v\n", err)
					os.Exit(1)
				}
				internalToken = strings.TrimSpace(string(byteToken))
			}
			fmt.Println("Exchanging local operator credential for CLI session...")
			authResp, err := apiClient.ExchangeLocalCLISession(internalToken)
			if err != nil {
				fmt.Printf("❌ Local session exchange failed: %v\n", err)
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
			fmt.Printf("✅ Successfully authenticated local session as %s\n", authResp.User.Email)
			return
		}

		if token != "" {
			fmt.Println("Verifying token...")
			user, err := apiClient.VerifyToken(token)
			if err != nil {
				fmt.Printf("❌ Token verification failed: %v\n", err)
				os.Exit(1)
			}
			if err := config.Save(&config.Config{
				ServerURL: serverURL,
				Token:     token,
				Email:     user.Email,
			}); err != nil {
				fmt.Printf("❌ Failed to save configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Successfully authenticated via token as %s\n", user.Email)
			return
		}

		email := strings.TrimSpace(loginEmail)
		password := loginPassword

		if email == "" {
			fmt.Print("Email (or leave blank to paste a token): ")
			input, _ := reader.ReadString('\n')
			email = strings.TrimSpace(input)
		}

		if email == "" {
			fmt.Print("Personal Access Token / API Token: ")
			byteToken, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				fmt.Printf("❌ Failed to read token: %v\n", err)
				os.Exit(1)
			}
			pastedToken := strings.TrimSpace(string(byteToken))
			if pastedToken == "" {
				fmt.Println("Error: email or token is required.")
				os.Exit(1)
			}
			fmt.Println("Verifying token...")
			user, err := apiClient.VerifyToken(pastedToken)
			if err != nil {
				fmt.Printf("❌ Token verification failed: %v\n", err)
				os.Exit(1)
			}
			if err := config.Save(&config.Config{
				ServerURL: serverURL,
				Token:     pastedToken,
				Email:     user.Email,
			}); err != nil {
				fmt.Printf("❌ Failed to save configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Successfully authenticated via token as %s\n", user.Email)
			return
		}

		if password == "" {
			fmt.Print("Password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				fmt.Println("Error reading password:", err)
				os.Exit(1)
			}
			password = string(bytePassword)
		}

		fmt.Println("Authenticating...")

		authResp, err := apiClient.Login(email, password)
		if err != nil && strings.Contains(err.Error(), "2FA code required") {
			fmt.Print("Enter 2FA Code: ")
			byteTotp, readErr := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if readErr != nil {
				fmt.Printf("❌ Failed to read 2FA code: %v\n", readErr)
				os.Exit(1)
			}
			authResp, err = apiClient.LoginWithTOTP(email, password, strings.TrimSpace(string(byteTotp)))
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
	loginCmd.Flags().StringVarP(&loginServerURL, "server", "s", "", "Codedock server API URL")
	loginCmd.Flags().StringVarP(&loginToken, "token", "t", "", "Personal access token or bearer token")
	loginCmd.Flags().StringVar(&loginTokenFile, "token-file", "", "File containing a personal access token")
	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "Account email address")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Account password")
	loginCmd.Flags().BoolVar(&loginLocal, "local", false, "Bootstrap CLI session using local installation operator token")
	rootCmd.AddCommand(loginCmd)
}
