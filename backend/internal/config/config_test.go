package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("ICTSC_API_ADDRESS", "")
	t.Setenv("ICTSC_API_READ_HEADER_TIMEOUT", "")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Address != ":8080" {
		t.Errorf("Address = %q, want %q", got.Address, ":8080")
	}
	if got.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout = %s, want %s", got.ReadHeaderTimeout, 5*time.Second)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("ICTSC_API_ADDRESS", ":9090")
	t.Setenv("ICTSC_API_READ_HEADER_TIMEOUT", "3s")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Address != ":9090" {
		t.Errorf("Address = %q, want %q", got.Address, ":9090")
	}
	if got.ReadHeaderTimeout != 3*time.Second {
		t.Errorf("ReadHeaderTimeout = %s, want %s", got.ReadHeaderTimeout, 3*time.Second)
	}
}

func TestLoadRejectsInvalidTimeout(t *testing.T) {
	t.Setenv("ICTSC_API_READ_HEADER_TIMEOUT", "invalid")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
}
