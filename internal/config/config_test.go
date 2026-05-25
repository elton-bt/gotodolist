package config

import (
	"net"
	"testing"
)

func TestVersionUsesRuntimeEnvWhenBuildVersionIsDefault(t *testing.T) {
	originalVersion := version
	version = "dev"
	defer func() {
		version = originalVersion
	}()

	t.Setenv("APP_VERSION", "1.2.0")

	if got := Version(); got != "1.2.0" {
		t.Fatalf("expected runtime APP_VERSION to be used, got %q", got)
	}
}

func TestVersionPrefersEmbeddedBuildVersion(t *testing.T) {
	originalVersion := version
	version = "2.3.4"
	defer func() {
		version = originalVersion
	}()

	t.Setenv("APP_VERSION", "1.2.0")

	if got := Version(); got != "2.3.4" {
		t.Fatalf("expected embedded build version to win, got %q", got)
	}
}

func TestVersionFallsBackToDevWhenUnset(t *testing.T) {
	originalVersion := version
	version = ""
	defer func() {
		version = originalVersion
	}()

	t.Setenv("APP_VERSION", "")

	if got := Version(); got != "dev" {
		t.Fatalf("expected dev fallback, got %q", got)
	}
}

func TestFirstNonLoopbackIPPrefersIPv4(t *testing.T) {
	addrs := []net.Addr{
		&net.IPNet{IP: net.ParseIP("127.0.0.1")},
		&net.IPNet{IP: net.ParseIP("fd00::10")},
		&net.IPNet{IP: net.ParseIP("10.42.0.15")},
	}

	if got := firstNonLoopbackIP(addrs); got != "10.42.0.15" {
		t.Fatalf("expected IPv4 address to be preferred, got %q", got)
	}
}

func TestFirstNonLoopbackIPFallsBackToIPv6(t *testing.T) {
	addrs := []net.Addr{
		&net.IPNet{IP: net.ParseIP("127.0.0.1")},
		&net.IPNet{IP: net.ParseIP("fd00::10")},
	}

	if got := firstNonLoopbackIP(addrs); got != "fd00::10" {
		t.Fatalf("expected IPv6 fallback, got %q", got)
	}
}
