package redis

import (
	"context"
	"testing"

	goredis "github.com/redis/go-redis/v9"
)

func TestNewReturnsNilWhenAddrEmpty(t *testing.T) {
	client, err := New(context.Background(), Config{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if client != nil {
		t.Fatalf("New() client = %v, want nil", client)
	}
}

func TestCloseAllowsNilClient(t *testing.T) {
	if err := Close(nil); err != nil {
		t.Fatalf("Close(nil) error = %v", err)
	}
}

func TestCloseClosesClient(t *testing.T) {
	client := goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:0"})

	if err := Close(client); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := client.Ping(context.Background()).Err(); err == nil {
		t.Fatal("Ping() error = nil, want closed client error")
	}
}

func TestKeyBuildsColonSeparatedRedisKey(t *testing.T) {
	got := Key("eugl:restaurant:pos:catalog:v1", int64(1), int64(10), "dine_in")
	want := "eugl:restaurant:pos:catalog:v1:1:10:dine_in"
	if got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}
}

func TestKeyPrefixAddsTrailingSeparator(t *testing.T) {
	got := KeyPrefix("eugl:restaurant:pos:catalog:v1", int64(1), int64(10))
	want := "eugl:restaurant:pos:catalog:v1:1:10:"
	if got != want {
		t.Fatalf("KeyPrefix() = %q, want %q", got, want)
	}
}
