package lushauth

import (
	"fmt"
	"strings"
)

// Permitter defines the behavior of allowing access.
type Permitter interface {
	Permit(c Consumer) error
}

// UserPolicy defines what users to grant access for.
type UserPolicy []string

// Permit a consumer or return an error.
func (p UserPolicy) Permit(c Consumer) error {
	if !c.HasAnyUUID(p...) {
		return p
	}
	return nil
}

func (p UserPolicy) Error() string {
	return fmt.Sprintf("need to be a specific user")
}

// RolePolicy defines what roles to grant access for.
type RolePolicy []string

// Permit a consumer or return an error.
func (p RolePolicy) Permit(c Consumer) error {
	if !c.HasAnyRole(p...) {
		return p
	}
	return nil
}

func (p RolePolicy) Error() string {
	return fmt.Sprintf("need to have any of the %s roles", strings.Join(quote(p...), ", "))
}

// GrantPolicy defines what grants required for access.
type GrantPolicy []string

// Permit a consumer or return an error.
func (p GrantPolicy) Permit(c Consumer) error {
	if !c.HasAnyGrant(p...) {
		return p
	}
	return nil
}

func (p GrantPolicy) Error() string {
	return fmt.Sprintf("need to have any of the %s grants", strings.Join(quote(p...), ", "))
}

// MarketPolicy defines what roles to allow access for in a given market.
type MarketPolicy struct {
	ID    string
	Roles []string
}

// Permit a consumer or return an error.
func (p MarketPolicy) Permit(c Consumer) error {
	if !c.HasAnyMarketRole(p.ID, p.Roles...) {
		return p
	}
	return nil
}

func (p MarketPolicy) Error() string {
	return fmt.Sprintf("need to be a part of the %q market with any of the %s market roles", p.ID, strings.Join(quote(p.Roles...), ", "))
}

// AnyPolicy defines a policy made up of multiple other policies where any of them will permit access.
type AnyPolicy []Permitter

// Permit a consumer or return an error.
func (p AnyPolicy) Permit(c Consumer) error {
	var errs []error
	for _, policy := range p {
		if err := policy.Permit(c); err != nil {
			errs = append(errs, err)
			continue
		}
		return nil
	}
	if len(errs) >= 1 {
		return errs[0]
	}
	return nil
}

// AllPolicy defines a policy made up of multiple other policies where all of them are required for access to be permitted.
type AllPolicy []Permitter

// Permit a consumer or return an error.
func (p AllPolicy) Permit(c Consumer) error {
	for _, policy := range p {
		if err := policy.Permit(c); err != nil {
			return err
		}
	}
	return nil
}

// quote will take a slice of strings and quote each of them
func quote(unquoted ...string) []string {
	quoted := make([]string, len(unquoted))
	for i, s := range unquoted {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return quoted
}
