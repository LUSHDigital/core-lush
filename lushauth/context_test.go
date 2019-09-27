package lushauth_test

import (
	"context"
	"corelush/lushauth"
	"testing"

	"github.com/LUSHDigital/core/test"
)

var (
	ctx context.Context
)

func ExampleContextWithConsumer() {
	ctx = lushauth.ContextWithConsumer(context.Background(), lushauth.Consumer{
		ID:     999,
		Grants: []string{"foo"},
	})
}

func ExampleConsumerFromContext() {
	consumer := lushauth.ConsumerFromContext(ctx)
	consumer.IsUser(999)
}

func TestContext(t *testing.T) {
	ctx = lushauth.ContextWithConsumer(context.Background(), lushauth.Consumer{
		ID:     999,
		Grants: []string{"foo"},
	})
	consumer := lushauth.ConsumerFromContext(ctx)
	test.Equals(t, true, consumer.IsUser(999))

	lushauth.ConsumerFromContext(context.Background())
	test.Equals(t, false, consumer.IsUser(0))
}
