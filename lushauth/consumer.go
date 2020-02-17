package lushauth

// Consumer represents an API user for the LUSH infrastructure.
type Consumer struct {
	// ID is a unique identifier for a user but should not be used in favour of UUID.
	ID int64 `json:"id"`
	// UUID is the unique identifier for a user.
	UUID string `json:"uuid"`
	// FirstName is the given name of a user.
	FirstName string `json:"first_name"`
	// LastName is the surname of a user.
	LastName string `json:"last_name"`
	// Language is the preferred language of a user.
	Language string `json:"language"`
	// Grants are any specific, given permissions for a user.
	// e.g. products.create, pages.read or tills.close
	Grants []string `json:"grants"`
	// Roles are what purpose a user server within the context of LUSH
	// e.g. guest, staff, creator or admin
	Roles []string `json:"roles"`
	// Needs are things that the user needs to do and that a front-end can react to.
	// e.g. password_reset, confirm_email or accept_terms
	Needs []string `json:"needs"`
	// Markets the user belongs to.
	// e.g. "gb", "de", etc...
	Markets []Market `json:"markets"`
}

// Market represents a market attached to an API user.
type Market struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

// HasAnyGrant checks if a consumer possess any of a given set of grants
func (c *Consumer) HasAnyGrant(grants ...string) bool {
	return hasAny(c.Grants, grants...)
}

// HasNoMatchingGrant checks if a consumer is missing any of a given set of grants
func (c Consumer) HasNoMatchingGrant(grants ...string) bool {
	return hasNoMatching(c.Grants, grants...)
}

// HasAnyRole checks if a consumer possess any of a given set of roles
func (c *Consumer) HasAnyRole(roles ...string) bool {
	return hasAny(c.Roles, roles...)
}

// HasNoMatchingRole checks if a consumer is missing any of a given set of roles
func (c *Consumer) HasNoMatchingRole(roles ...string) bool {
	return hasNoMatching(c.Roles, roles...)
}

// HasAnyNeed checks if a consumer has any of the given needs
func (c *Consumer) HasAnyNeed(needs ...string) bool {
	return hasAny(c.Needs, needs...)
}

// HasNoMatchingNeed checks if a consumer has any of the given needs
func (c *Consumer) HasNoMatchingNeed(needs ...string) bool {
	return hasNoMatching(c.Needs, needs...)
}

// IsUser checks if a consumer has the same ID as a user
func (c *Consumer) IsUser(id int64) bool {
	return c.ID == id
}

// HasAnyUUID checks if a consumer has the same uuid as a user
func (c *Consumer) HasAnyUUID(ids ...string) bool {
	return hasAny([]string{c.UUID}, ids...)
}

// HasAnyMarketRole checks if a user has any role in a given market.
func (c Consumer) HasAnyMarketRole(id string, roles ...string) bool {
	for _, m := range c.Markets {
		if m.ID == id {
			return hasAny(m.Roles, roles...)
		}
	}
	return false
}

// hasAny checks if a set contains any of the given members.
func hasAny(set []string, members ...string) bool {
	for _, member := range members {
		for _, m := range set {
			if member == m {
				return true
			}
		}
	}
	return false
}

// hasNoMatching checks if a set does not contain any and all the given members.
func hasNoMatching(set []string, members ...string) bool {
	for _, member := range members {
		for _, m := range set {
			if member == m {
				return false
			}
		}
	}
	return true
}
