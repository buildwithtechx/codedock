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
		if strings.ContainsAny(domain, "/: \t") || !strings.Contains(domain, ".") {
			exitError("Enter a domain name without a scheme, port, or path")
		}
		tlsEmail = strings.TrimSpace(promptOptional(fmt.Sprintf("Certificate email [%s]: ", email)))
		if tlsEmail == "" {
			tlsEmail = email
		}
		if _, err := mail.ParseAddress(tlsEmail); err != nil {
			exitError("Invalid certificate email: %v", err)
		}
		if err := saveSetupOptions(db, vault, domain, tlsEmail); err != nil {
			exitError("Save setup settings: %v", err)
		}
		fmt.Printf("Domain saved. Point *.%s at your server and restart Codedock to enable HTTPS.\n", domain)
	}
	fmt.Printf("Dashboard: http://localhost:%d\n", config.Get().Server.Port)
	fmt.Println("Optional email and OAuth providers can be configured in dashboard settings.")
}

func createSetupOwner(repo *repositories.UserRepo) string {
	fmt.Println("Create your instance owner account.")
	email := prompt("Email: ")
	if _, err := mail.ParseAddress(email); err != nil {
		exitError("Invalid email: %v", err)
	}
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
	if err := config.SaveSelfHostedOptions(config.Get(), domain, tlsEmail); err != nil {
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
