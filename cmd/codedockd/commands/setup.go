package commands

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

func runSetup() {
	fmt.Println("Codedock setup")
	_, db, vault := InitDataDir()
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Close setup database: %v\n", err)
		}
	}()
	if config.Get().Cloud.Enabled {
		exitError("The setup wizard is for self-hosted installations")
	}
	repo := repositories.NewUserRepo(db)
	users, _, err := repo.ListUsers(context.Background(), 1, 0)
	if err != nil {
		exitError("Check existing accounts: %v", err)
	}
	email := ""
	if len(users) == 0 {
		email = createSetupOwner(repo)
	} else {
		email = users[0].Email
		fmt.Println("An account already exists. Use 'codedockd reset-password' to reset its password.")
	}
	domain := strings.TrimPrefix(strings.TrimSpace(promptOptional("Optional domain for apps (e.g. apps.yourdomain.com; Enter to use automatic DNS): ")), "*.")
	tlsEmail := ""
	if domain != "" {
		if err := validateSetupDomain(domain); err != nil {
			exitError("Invalid app domain: %v", err)
		}
		tlsEmail = strings.TrimSpace(promptOptional(fmt.Sprintf("Certificate email [%s]: ", email)))
		if tlsEmail == "" {
			tlsEmail = email
		}
		address, err := mail.ParseAddress(tlsEmail)
		if err != nil {
			exitError("Invalid certificate email: %v", err)
		}
		tlsEmail = address.Address
		fmt.Printf("Domain saved. Point *.%s at your server and restart Codedock to enable HTTPS.\n", domain)
	}
	if err := saveSetupOptions(db, vault, domain, tlsEmail); err != nil {
		exitError("Save setup settings: %v", err)
	}
	if domain == "" {
		fmt.Println("Automatic DNS selected. Restart Codedock to apply routing settings.")
	}
	fmt.Printf("Dashboard: http://localhost:%d\n", config.Get().Server.Port)
	fmt.Println("Optional email and OAuth providers can be configured in dashboard settings.")
}

func createSetupOwner(repo *repositories.UserRepo) string {
	fmt.Println("Create your instance owner account.")
	email := prompt("Email: ")
	address, err := mail.ParseAddress(email)
	if err != nil {
		exitError("Invalid email: %v", err)
	}
	email = address.Address
	name := strings.TrimSpace(prompt("Name: "))
	if name == "" {
		exitError("Name is required")
	}
	password := promptPassword("Password: ")
	if password != promptPassword("Confirm password: ") {
		exitError("Passwords do not match")
	}
	if err := utils.ValidatePassword(password); err != nil {
		exitError("Invalid password: %v", err)
	}
	hashed, err := utils.HashPassword(password)
	if err != nil {
		exitError("Hash password: %v", err)
	}
	user := &models.User{ID: uuid.New().String(), Email: email, Name: name, PasswordHash: hashed, Role: models.UserRoleOwner, IsActive: true, EmailVerified: true}
	if err := repo.CreateUser(context.Background(), user); err != nil {
		exitError("Create instance owner: %v", err)
	}
	fmt.Printf("Owner account created for %s.\n", email)
	return email
}

func saveSetupOptions(db *sql.DB, vault *utils.Vault, domain, tlsEmail string) error {
	if err := config.SaveSelfHostedOptions(context.Background(), config.Get(), repositories.NewSelfHostedRepo(db), domain, tlsEmail); err != nil {
		return err
	}
	repo := repositories.NewSettingsRepo(db, vault)
	settings, err := repo.GetServerSettings(context.Background())
	if err != nil {
		return fmt.Errorf("load server settings: %w", err)
	}
	settings.DefaultWildcardDomain = domain
	if err := repo.UpdateServerSettings(context.Background(), settings); err != nil {
		return fmt.Errorf("save wildcard domain: %w", err)
	}
	return nil
}

func exitError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func validateSetupDomain(domain string) error {
	if len(domain) > 253 || !strings.Contains(domain, ".") {
		return fmt.Errorf("use a fully qualified domain name")
	}
	for _, label := range strings.Split(domain, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return fmt.Errorf("invalid DNS label")
		}
		for _, character := range label {
			if !(character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' || character >= '0' && character <= '9' || character == '-') {
				return fmt.Errorf("only letters, digits, and hyphens are allowed")
			}
		}
	}
	return nil
}
