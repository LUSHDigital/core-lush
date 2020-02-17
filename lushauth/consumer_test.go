package lushauth_test

import (
	"fmt"
	"testing"

	"github.com/LUSHDigital/core-lush/lushauth"

	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/uuid"
)

type hasAnyTest struct {
	name     string
	values   []string
	expected bool
}

func hasAnyCases(kind string) []hasAnyTest {
	return []hasAnyTest{
		{
			name:     fmt.Sprintf("when using one %s that exists", kind),
			values:   []string{"test.foo"},
			expected: true,
		},
		{
			name:     fmt.Sprintf("when using two %ss where one does not exist", kind),
			values:   []string{"test.foo", "doesnot.exist"},
			expected: true,
		},
		{
			name:     fmt.Sprintf("when using one %s that does not exist", kind),
			values:   []string{"doesnot.exist"},
			expected: false,
		},
		{
			name:     fmt.Sprintf("when using two %ss that does not exist", kind),
			values:   []string{"doesnot.exist", "has.no.access"},
			expected: false,
		},
	}
}

func TestConsumer_HasAnyGrant(t *testing.T) {
	for _, c := range hasAnyCases("grant") {
		test.Equals(t, c.expected, consumer.HasAnyGrant(c.values...))
	}
}

func TestConsumer_HasNoMatchingGrant(t *testing.T) {
	for _, c := range hasAnyCases("grant") {
		test.Equals(t, !c.expected, consumer.HasNoMatchingGrant(c.values...))
	}
}

func TestConsumer_HasAnyNeed(t *testing.T) {
	for _, c := range hasAnyCases("need") {
		test.Equals(t, c.expected, consumer.HasAnyNeed(c.values...))
	}
}

func TestConsumer_HasNoMatchingNeed(t *testing.T) {
	for _, c := range hasAnyCases("need") {
		test.Equals(t, !c.expected, consumer.HasNoMatchingNeed(c.values...))
	}
}

func TestConsumer_HasAnyRole(t *testing.T) {
	for _, c := range hasAnyCases("role") {
		test.Equals(t, c.expected, consumer.HasAnyRole(c.values...))
	}
}

func TestConsumer_HasNoMatchingRole(t *testing.T) {
	for _, c := range hasAnyCases("role") {
		test.Equals(t, !c.expected, consumer.HasNoMatchingRole(c.values...))
	}
}

func TestConsumer_IsUser(t *testing.T) {
	consumer := &lushauth.Consumer{
		ID: 1,
	}
	t.Run("when its the same user", func(t *testing.T) {
		test.Equals(t, true, consumer.IsUser(1))
	})
	t.Run("when its not the same user", func(t *testing.T) {
		test.Equals(t, false, consumer.IsUser(2))
	})
}

func TestConsumer_HasUUID(t *testing.T) {
	id1 := uuid.Must(uuid.NewV4()).String()
	id2 := uuid.Must(uuid.NewV4()).String()
	consumer := &lushauth.Consumer{
		UUID: id1,
	}
	t.Run("when its the same user", func(t *testing.T) {
		test.Equals(t, true, consumer.HasAnyUUID(id1))
	})
	t.Run("when its not the same user", func(t *testing.T) {
		test.Equals(t, false, consumer.HasAnyUUID(id2))
	})
}

func TestConsumer_HasAnyMarketRole(t *testing.T) {
	consumer := &lushauth.Consumer{
		Markets: []lushauth.Market{
			{
				ID:    "gb",
				Roles: []string{"one", "two"},
			},
			{
				ID:    "fr",
				Roles: []string{"three", "four"},
			},
		},
	}
	t.Run("when user has role in 'gb' market", func(t *testing.T) {
		test.Equals(t, true, consumer.HasAnyMarketRole("gb", "one"))
	})
	t.Run("when user doesn't have role in 'fr' market", func(t *testing.T) {
		test.Equals(t, false, consumer.HasAnyMarketRole("fr", "one"))
	})
	t.Run("when user has role in 'fr' market", func(t *testing.T) {
		test.Equals(t, true, consumer.HasAnyMarketRole("fr", "four", "five"))
	})
	t.Run("when user doesn't belong to the market at all", func(t *testing.T) {
		test.Equals(t, false, consumer.HasAnyMarketRole("se", "two"))
	})
}
