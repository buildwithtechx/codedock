package http

import (
	"context"

	"codedock/internal/modules/billing"
	"codedock/internal/modules/users"
)

type billingUserAdapter struct {
	repo users.UserRepository
}

func (a *billingUserAdapter) GetUserByID(ctx context.Context, id string) (*billing.User, error) {
	if a.repo == nil {
		return nil, nil
	}
	u, err := a.repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return toBillingUser(u), nil
}

func (a *billingUserAdapter) GetUserByStripeCustomerID(ctx context.Context, customerID string) (*billing.User, error) {
	if a.repo == nil {
		return nil, nil
	}
	u, err := a.repo.GetUserByStripeCustomerID(ctx, customerID)
	if err != nil || u == nil {
		return nil, err
	}
	return toBillingUser(u), nil
}

func (a *billingUserAdapter) UpdateUser(ctx context.Context, user *billing.User) error {
	if a.repo == nil || user == nil {
		return nil
	}
	u, err := a.repo.GetUserByID(ctx, user.ID)
	if err != nil || u == nil {
		return err
	}
	u.StripeCustomerID = user.StripeCustomerID
	u.StripeSubscriptionID = user.StripeSubscriptionID
	u.StripePriceID = user.StripePriceID
	u.PlanType = user.PlanType
	return a.repo.UpdateUser(ctx, u)
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
