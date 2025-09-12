package license

import (
	"strings"
	"testing"
)

func TestGetMACAddress(t *testing.T) {
	mac := GetMACAddress()

	if mac != "" {
		if len(mac) != 12 {
			t.Errorf("MAC address should be 12 characters long, got %d", len(mac))
		}

		for _, char := range mac {
			if !strings.Contains("0123456789ABCDEF", string(char)) {
				t.Errorf("MAC address should contain only hex characters, got %c", char)
			}
		}
	}
}

func TestGetIPAddress(t *testing.T) {
	ip := GetIPAddress()

	if ip != "" {
		if !strings.Contains(ip, ".") {
			t.Errorf("IP address should contain dots, got %s", ip)
		}
	}
}

func TestGetSystemFingerprint(t *testing.T) {
	fingerprint := GetSystemFingerprint()

	if len(fingerprint) != 64 {
		t.Errorf("Fingerprint should be 64 characters long (SHA256), got %d", len(fingerprint))
	}

	for _, char := range fingerprint {
		if !strings.Contains("0123456789abcdef", string(char)) {
			t.Errorf("Fingerprint should contain only lowercase hex characters, got %c", char)
		}
	}
}

func TestGetMultiHostFingerprint(t *testing.T) {
	fingerprint := GetMultiHostFingerprint()

	if len(fingerprint) != 64 {
		t.Errorf("Multi-host fingerprint should be 64 characters long (SHA256), got %d", len(fingerprint))
	}

	for _, char := range fingerprint {
		if !strings.Contains("0123456789abcdef", string(char)) {
			t.Errorf("Multi-host fingerprint should contain only lowercase hex characters, got %c", char)
		}
	}
}
