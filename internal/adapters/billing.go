package adapters

import (
	"context"

	"codedock/internal/modules/billing"
	"codedock/internal/modules/users"
)

type BillingUserAdapter struct {
	Repo users.UserRepository
}

func NewBillingUserAdapter(repo users.UserRepository) *BillingUserAdapter {
	return &BillingUserAdapter{Repo: repo}
}

func (a *BillingUserAdapter) GetUserByID(ctx context.Context, id string) (*billing.User, error) {
	if a.Repo == nil {
		return nil, nil
	}
	u, err := a.Repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return toBillingUser(u), nil
}

func (a *BillingUserAdapter) GetUserByStripeCustomerID(ctx context.Context, customerID string) (*billing.User, error) {
	if a.Repo == nil {
		return nil, nil
	}
	u, err := a.Repo.GetUserByStripeCustomerID(ctx, customerID)
	if err != nil || u == nil {
		return nil, err
	}
	return toBillingUser(u), nil
}

func (a *BillingUserAdapter) UpdateUser(ctx context.Context, user *billing.User) error {
	if a.Repo == nil || user == nil {
		return nil
	}
	u, err := a.Repo.GetUserByID(ctx, user.ID)
	if err != nil || u == nil {
		return err
	}
	u.StripeCustomerID = user.StripeCustomerID
	u.StripeSubscriptionID = user.StripeSubscriptionID
	u.StripePriceID = user.StripePriceID
	u.PlanType = user.PlanType
	return a.Repo.UpdateUser(ctx, u)
}

func toBillingUser(u *users.User) *billing.User {
	if u == nil {
		return nil
	}
	return &billing.User{
		ID:                   u.ID,
		Email:                u.Email,
		Name:                 u.Name,
		Role:                 string(u.Role),
		PlanType:             u.PlanType,
		StripeCustomerID:     u.StripeCustomerID,
		StripeSubscriptionID: u.StripeSubscriptionID,
		StripePriceID:        u.StripePriceID,
	}
}
