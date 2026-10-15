package service

import "testing"

func TestValidOAuthNext(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{value: "/", valid: true},
		{value: "/admin/problems/ABC", valid: true},
		{value: "https://evil.example", valid: false},
		{value: "//evil.example", valid: false},
		{value: "/\\\\evil.example", valid: false},
		{value: "/%5C%5Cevil.example", valid: false},
		{value: "/%255C%255Cevil.example", valid: false},
		{value: "/admin/?next=evil", valid: false},
		{value: "/admin/#fragment", valid: false},
		{value: "/admin/\nlocation:https://evil.example", valid: false},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			if got := validOAuthNext(test.value); got != test.valid {
				t.Fatalf("validOAuthNext(%q) = %v, want %v", test.value, got, test.valid)
			}
		})
	}
}
