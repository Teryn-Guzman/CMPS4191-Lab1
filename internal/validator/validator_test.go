package validator

import "testing"

func TestValidator(t *testing.T) {
	v := New()
	v.Check(false, "name", "must be provided")
	v.Check(true, "email", "must be provided")

	if v.Valid() {
		t.Fatal("expected validator to be invalid")
	}
	if len(v.Errors["name"]) != 1 {
		t.Fatalf("expected one name error, got %d", len(v.Errors["name"]))
	}
}

func TestValidateContactEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{name: "valid value", email: "user@example.com", valid: true},
		{name: "missing value", email: "", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			ValidateContactEmail(v, tt.email)
			if v.Valid() != tt.valid {
				t.Fatalf("expected valid=%t, got %t", tt.valid, v.Valid())
			}
		})
	}
}
