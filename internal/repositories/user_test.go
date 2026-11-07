package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
)

func TestUserRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewUserRepo(db)
	stripeID := "cus_test123"
	user := &models.User{
		Email:            "ada@example.com",
		Name:             "Ada",
		PasswordHash:     "hash",
		Role:             models.UserRoleOwner,
		IsActive:         true,
		EmailVerified:    true,
		PlanType:         "pro",
		StripeCustomerID: &stripeID,
	}
	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.ID == "" {
		t.Fatal("expected generated user id")
	}
	byID, err := repo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if byID.Email != user.Email || !byID.IsActive || byID.PlanType != "pro" {
		t.Fatalf("unexpected user row: %+v", byID)
	}
	if byID.StripeCustomerID == nil || *byID.StripeCustomerID != stripeID {
		t.Fatal("expected stripe customer id roundtrip")
	}
	if byID.CreatedAt.IsZero() || byID.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to roundtrip")
	}
	byEmail, err := repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if byEmail.ID != user.ID {
		t.Fatal("email lookup returned wrong user")
	}
	byStripe, err := repo.GetUserByStripeCustomerID(ctx, stripeID)
	if err != nil {
		t.Fatalf("get by stripe id: %v", err)
	}
	if byStripe.ID != user.ID {
		t.Fatal("stripe lookup returned wrong user")
	}
	count, err := repo.CountUsers(ctx)
	if err != nil {
		t.Fatalf("count users: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 user, got %d", count)
	}
	users, total, err := repo.ListUsers(ctx, 10, 0)
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if total != 1 || len(users) != 1 || users[0].ID != user.ID {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(users))
	}
	login := time.Now().UTC().Truncate(time.Second)
	user.Name = "Ada Lovelace"
	user.LastLogin = &login
	if err := repo.UpdateUser(ctx, user); err != nil {
		t.Fatalf("update user: %v", err)
	}
	updated, err := repo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get updated user: %v", err)
	}
	if updated.Name != "Ada Lovelace" || updated.LastLogin == nil {
		t.Fatalf("update did not persist: %+v", updated)
	}
	secret, codes, err := repo.GetUserTOTPSecret(ctx, user.ID)
	if err != nil {
		t.Fatalf("get totp secret: %v", err)
	}
	if secret != "" || len(codes) != 0 {
		t.Fatalf("expected empty totp material, got %q %v", secret, codes)
	}
	expiry := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	pat := &models.PersonalAccessToken{
		UserID:      user.ID,
		Name:        "ci",
		TokenHash:   "hash",
		Prefix:      "cdk_",
		AccessLevel: "read",
		ExpiresAt:   &expiry,
	}
	if err := repo.CreatePAT(ctx, pat); err != nil {
		t.Fatalf("create pat: %v", err)
	}
	pats, err := repo.ListPATs(ctx, user.ID)
	if err != nil {
		t.Fatalf("list pats: %v", err)
	}
	if len(pats) != 1 || pats[0].ID != pat.ID {
		t.Fatalf("expected 1 pat, got %d", len(pats))
	}
	if pats[0].ExpiresAt == nil || pats[0].CreatedAt.IsZero() {
		t.Fatal("expected pat timestamps to roundtrip")
	}
	if err := repo.DeletePAT(ctx, pat.ID, user.ID); err != nil {
		t.Fatalf("delete pat: %v", err)
	}
	if err := repo.DeleteUser(ctx, user.ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if _, err := repo.GetUserByID(ctx, user.ID); err == nil {
		t.Fatal("expected deleted user to be gone")
	}
}
