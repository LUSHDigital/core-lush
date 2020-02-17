package lushauth_test

import (
	"testing"

	"github.com/LUSHDigital/core-lush/lushauth"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/uuid"
)

var (
	GuestID      = uuid.Must(uuid.NewV4()).String()
	PolicyUserID = uuid.Must(uuid.NewV4()).String()

	Guest = lushauth.Consumer{
		UUID:    GuestID,
		Roles:   []string{},
		Markets: []lushauth.Market{},
	}
	Staff = lushauth.Consumer{
		UUID:  PolicyUserID,
		Roles: []string{"staff"},
		Markets: []lushauth.Market{
			{
				ID:    "gb",
				Roles: []string{"staff"},
			},
		},
	}
	Manager = lushauth.Consumer{
		UUID:  PolicyUserID,
		Roles: []string{"staff"},
		Markets: []lushauth.Market{
			{
				ID:    "gb",
				Roles: []string{"manager"},
			},
			{
				ID:    "se",
				Roles: []string{"manager"},
			},
		},
	}
	Admin = lushauth.Consumer{
		UUID:   PolicyUserID,
		Roles:  []string{"admin"},
		Grants: []string{"users.delete", "users.create"},
		Markets: []lushauth.Market{
			{
				ID:    "gb",
				Roles: []string{"manager"},
			},
		},
	}
	Deleter = lushauth.Consumer{
		UUID:   PolicyUserID,
		Grants: []string{"users.delete"},
	}
)

func ExampleRolePolicy() {
	policy := lushauth.RolePolicy{"admin", "staff"}
	policy.Permit(consumer)
}

func TestRolePolicy_Permit(t *testing.T) {
	var (
		StaffPolicy = lushauth.RolePolicy{"admin", "staff"}
		AdminPolicy = lushauth.RolePolicy{"admin"}
	)
	type Test struct {
		name     string
		consumer lushauth.Consumer
		policy   lushauth.Permitter
		expected error
	}
	cases := []Test{
		{
			name:     "staff policy with permitted consumer",
			consumer: Staff,
			policy:   StaffPolicy,
			expected: nil,
		},
		{
			name:     "staff policy with consumer lacking staff role",
			consumer: Guest,
			policy:   StaffPolicy,
			expected: StaffPolicy,
		},
		{
			name:     "admin policy with permitted consumer",
			consumer: Admin,
			policy:   AdminPolicy,
			expected: nil,
		},
		{
			name:     "admin policy with consumer lacking admin role",
			consumer: Staff,
			policy:   AdminPolicy,
			expected: AdminPolicy,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AssertPermit(t, c.expected, c.policy, c.consumer)
		})
	}
}

func ExampleGrantPolicy() {
	policy := lushauth.GrantPolicy{"users.delete"}
	policy.Permit(consumer)
}

func TestGrantPolicy_Permit(t *testing.T) {
	var (
		DeleteUsersPolicy = lushauth.GrantPolicy{"users.delete"}
	)
	type Test struct {
		name     string
		consumer lushauth.Consumer
		policy   lushauth.Permitter
		expected error
	}
	cases := []Test{
		{
			name:     "delete users policy with permitted consumer",
			consumer: Deleter,
			policy:   DeleteUsersPolicy,
			expected: nil,
		},
		{
			name:     "delete users policy with consumer lacking staff role",
			consumer: Guest,
			policy:   DeleteUsersPolicy,
			expected: DeleteUsersPolicy,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AssertPermit(t, c.expected, c.policy, c.consumer)
		})
	}
}

var UserID string

func ExampleUserPolicy() {
	policy := lushauth.UserPolicy{
		UserID, // UserID: "5d4b32f9-5954-41c3-a470-7d76317635a7"
	}
	policy.Permit(consumer)
}

func TestUserPolicy_Permit(t *testing.T) {
	var (
		SpecificUserPolicy = lushauth.UserPolicy{PolicyUserID}
	)
	type Test struct {
		name     string
		consumer lushauth.Consumer
		policy   lushauth.Permitter
		expected error
	}
	cases := []Test{
		{
			name:     "specific user policy with permitted consumer",
			consumer: Staff,
			policy:   SpecificUserPolicy,
			expected: nil,
		},
		{
			name:     "specific user policy with consumer not matching",
			consumer: Guest,
			policy:   SpecificUserPolicy,
			expected: SpecificUserPolicy,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AssertPermit(t, c.expected, c.policy, c.consumer)
		})
	}
}

func ExampleMarketPolicy() {
	policy := lushauth.MarketPolicy{
		ID: "gb",
		Roles: []string{
			"admin",
			"manager",
			"staff",
		},
	}
	policy.Permit(consumer)
}

func TestMarketPolicy_Permit(t *testing.T) {
	var (
		BritishMarketPolicy = lushauth.MarketPolicy{
			ID:    "gb",
			Roles: []string{"manager"},
		}
	)
	type Test struct {
		name     string
		consumer lushauth.Consumer
		policy   lushauth.Permitter
		expected error
	}
	cases := []Test{
		{
			name:     "british market with permitted consumer",
			consumer: Admin,
			policy:   BritishMarketPolicy,
			expected: nil,
		},
		{
			name:     "british market with consumer not part of the market",
			consumer: Guest,
			policy:   BritishMarketPolicy,
			expected: BritishMarketPolicy,
		},
		{
			name:     "british market with consumer part of the market but lacking the role",
			consumer: Guest,
			policy:   BritishMarketPolicy,
			expected: BritishMarketPolicy,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AssertPermit(t, c.expected, c.policy, c.consumer)
		})
	}
}

func ExampleAnyPolicy() {
	policy := lushauth.AnyPolicy{
		lushauth.GrantPolicy{"users.delete"},
		lushauth.RolePolicy{"admin"},
	}
	policy.Permit(consumer)
}

func TestAnyPolicy_Permit(t *testing.T) {
	var (
		GBPolicy = lushauth.MarketPolicy{
			ID:    "gb",
			Roles: []string{"staff"},
		}
		SEPolicy = lushauth.MarketPolicy{
			ID:    "se",
			Roles: []string{"staff"},
		}
		AnyMarketPolicy = lushauth.AnyPolicy{
			GBPolicy,
			SEPolicy,
		}
	)
	type Test struct {
		name     string
		consumer lushauth.Consumer
		policy   lushauth.Permitter
		expected error
	}
	cases := []Test{
		{
			name:     "any of no policies",
			consumer: Guest,
			policy:   lushauth.AnyPolicy{},
		},
		{
			name:     "any of two markets with permitted user",
			consumer: Staff,
			policy:   AnyMarketPolicy,
		},
		{
			name:     "any of two markets with user not part of any of the markets",
			consumer: Guest,
			policy:   AnyMarketPolicy,
			expected: GBPolicy,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AssertPermit(t, c.expected, c.policy, c.consumer)
		})
	}
}

func ExampleAllPolicy() {
	policy := lushauth.AllPolicy{
		lushauth.RolePolicy{"staff"},
		lushauth.MarketPolicy{
			ID:    "gb",
			Roles: []string{"manager"},
		},
	}
	policy.Permit(consumer)
}

func TestAllPolicy_Permit(t *testing.T) {
	var (
		GBPolicy = lushauth.MarketPolicy{
			ID:    "gb",
			Roles: []string{"staff", "manager"},
		}
		SEPolicy = lushauth.MarketPolicy{
			ID:    "se",
			Roles: []string{"staff", "manager"},
		}
		AllMarketsPolicy = lushauth.AllPolicy{
			GBPolicy,
			SEPolicy,
		}
	)
	type Test struct {
		name     string
		consumer lushauth.Consumer
		policy   lushauth.Permitter
		expected error
	}
	cases := []Test{
		{
			name:     "all of no policies",
			consumer: Guest,
			policy:   lushauth.AllPolicy{},
		},
		{
			name:     "all of two markets with permitted user",
			consumer: Manager,
			policy:   AllMarketsPolicy,
		},
		{
			name:     "all of two markets with user only part of one of the markets",
			consumer: Staff,
			policy:   AllMarketsPolicy,
			expected: SEPolicy,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AssertPermit(t, c.expected, c.policy, c.consumer)
		})
	}
}

func AssertPermit(t *testing.T, expected error, p lushauth.Permitter, c lushauth.Consumer) {
	t.Helper()
	err := p.Permit(c)
	if err != nil {
		t.Log(err)
	}
	test.Equals(t, expected, err)
}
