package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
)

func TestOrganizationRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	users := NewUserRepo(db)
	owner := &models.User{Email: "owner@example.com", Name: "Owner", PasswordHash: "hash", Role: models.UserRoleOwner, IsActive: true}
	if err := users.CreateUser(ctx, owner); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	memberUser := &models.User{Email: "member@example.com", Name: "Member", PasswordHash: "hash", Role: models.UserRoleMember, IsActive: true}
	if err := users.CreateUser(ctx, memberUser); err != nil {
		t.Fatalf("create member user: %v", err)
	}
	repo := NewOrganizationRepository(db)
	now := time.Now().UTC().Truncate(time.Second)
	org := &models.Organization{ID: uuid.NewString(), Name: "Acme", CreatedAt: now, UpdatedAt: now}
	if err := repo.Create(ctx, org); err != nil {
		t.Fatalf("create org: %v", err)
	}
	got, err := repo.GetByID(ctx, org.ID)
	if err != nil {
		t.Fatalf("get org: %v", err)
	}
	if got == nil || got.Name != "Acme" {
		t.Fatalf("unexpected org row: %+v", got)
	}
	if !got.CreatedAt.Truncate(time.Second).Equal(now) || !got.UpdatedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("org timestamps mismatch: %+v", got)
	}
	org.Name = "Acme Inc"
	if err := repo.Update(ctx, org); err != nil {
		t.Fatalf("update org: %v", err)
	}
	renamed, err := repo.GetByID(ctx, org.ID)
	if err != nil {
		t.Fatalf("get renamed org: %v", err)
	}
	if renamed == nil || renamed.Name != "Acme Inc" {
		t.Fatalf("update did not persist: %+v", renamed)
	}
	org2 := &models.Organization{ID: uuid.NewString(), Name: "Globex", CreatedAt: now, UpdatedAt: now}
	ownerMember := &models.OrganizationMember{
		ID:             uuid.NewString(),
		OrganizationID: org2.ID,
		UserID:         owner.ID,
		Email:          owner.Email,
		Permission:     models.MemberPermissionOwner,
		Status:         models.MemberStatusActive,
		InvitedAt:      now,
		AcceptedAt:     now,
	}
	if err := repo.CreateWithOwner(ctx, org2, ownerMember); err != nil {
		t.Fatalf("create org with owner: %v", err)
	}
	byUser, err := repo.ListByUser(ctx, owner.ID)
	if err != nil {
		t.Fatalf("list orgs by user: %v", err)
	}
	if len(byUser) != 1 || byUser[0].ID != org2.ID {
		t.Fatalf("expected 1 org for owner, got %d", len(byUser))
	}
	invite := &models.OrganizationMember{
		ID:             uuid.NewString(),
		OrganizationID: org.ID,
		UserID:         memberUser.ID,
		Email:          memberUser.Email,
		Permission:     models.MemberPermissionMember,
		Status:         models.MemberStatusPending,
		InvitedAt:      now,
	}
	if err := repo.AddMember(ctx, invite); err != nil {
		t.Fatalf("add member: %v", err)
	}
	fetched, err := repo.GetMember(ctx, org.ID, memberUser.ID)
	if err != nil {
		t.Fatalf("get member: %v", err)
	}
	if fetched == nil || fetched.ID != invite.ID || fetched.Status != models.MemberStatusPending {
		t.Fatalf("unexpected member row: %+v", fetched)
	}
	if !fetched.InvitedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("invited_at mismatch: %v", fetched.InvitedAt)
	}
	byID, err := repo.GetMemberByID(ctx, invite.ID)
	if err != nil {
		t.Fatalf("get member by id: %v", err)
	}
	if byID == nil || byID.Email != memberUser.Email {
		t.Fatalf("member by id mismatch: %+v", byID)
	}
	byEmail, err := repo.GetMemberByEmail(ctx, org.ID, memberUser.Email)
	if err != nil {
		t.Fatalf("get member by email: %v", err)
	}
	if byEmail == nil || byEmail.ID != invite.ID {
		t.Fatalf("member by email mismatch: %+v", byEmail)
	}
	members, err := repo.ListMembers(ctx, org.ID)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if len(members) != 1 || members[0].ID != invite.ID {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
	invites, err := repo.ListInvitesByEmail(ctx, memberUser.Email)
	if err != nil {
		t.Fatalf("list invites: %v", err)
	}
	if len(invites) != 1 || invites[0].ID != invite.ID {
		t.Fatalf("expected 1 invite, got %d", len(invites))
	}
	invite.Status = models.MemberStatusActive
	invite.AcceptedAt = now
	if err := repo.UpdateMember(ctx, invite); err != nil {
		t.Fatalf("update member: %v", err)
	}
	accepted, err := repo.GetMemberByID(ctx, invite.ID)
	if err != nil {
		t.Fatalf("get accepted member: %v", err)
	}
	if accepted == nil || accepted.Status != models.MemberStatusActive {
		t.Fatalf("member update did not persist: %+v", accepted)
	}
	if !accepted.AcceptedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("accepted_at mismatch: %v", accepted.AcceptedAt)
	}
	invites, err = repo.ListInvitesByEmail(ctx, memberUser.Email)
	if err != nil {
		t.Fatalf("list invites after accept: %v", err)
	}
	if len(invites) != 0 {
		t.Fatalf("expected 0 invites after accept, got %d", len(invites))
	}
	if err := repo.RemoveMember(ctx, invite.ID); err != nil {
		t.Fatalf("remove member: %v", err)
	}
	removed, err := repo.GetMemberByID(ctx, invite.ID)
	if err != nil {
		t.Fatalf("get removed member: %v", err)
	}
	if removed != nil {
		t.Fatal("expected removed member to be gone")
	}
	if err := repo.Delete(ctx, org.ID); err != nil {
		t.Fatalf("delete org: %v", err)
	}
	deleted, err := repo.GetByID(ctx, org.ID)
	if err != nil {
		t.Fatalf("get deleted org: %v", err)
	}
	if deleted != nil {
		t.Fatal("expected deleted org to be gone")
	}
	missing, err := repo.GetByID(ctx, "does-not-exist")
	if err != nil {
		t.Fatalf("get missing org: %v", err)
	}
	if missing != nil {
		t.Fatal("expected missing org lookup to return nil")
	}
}
