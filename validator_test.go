package validator

import (
	"testing"
)

type TestUser struct {
	Name     string
	Email    string
	Password string
	Age      int
	Role     string
	Tags     []string
}

func TestRequired(t *testing.T) {
	user := TestUser{Name: "", Email: "test@example.com"}
	v := New(&user)
	v.Required("Name")
	v.Required("Email")

	if v.IsValid() {
		t.Error("expected validation to fail for empty Name")
	}

	if _, exists := v.Errors["Name"]; !exists {
		t.Error("expected error for Name field")
	}

	if _, exists := v.Errors["Email"]; exists {
		t.Error("didn't expect error for Email field")
	}
}

func TestMinLength(t *testing.T) {
	user := TestUser{Password: "abc"}
	v := New(&user)
	v.MinLength("Password", 8)

	if v.IsValid() {
		t.Error("expected password to be too short")
	}

	user2 := TestUser{Password: "password123"}
	v2 := New(&user2)
	v2.MinLength("Password", 8)

	if !v2.IsValid() {
		t.Error("password should be valid")
	}
}

func TestMaxLength(t *testing.T) {
	user := TestUser{Name: "thisnameiswaytooolong"}
	v := New(&user)
	v.MaxLength("Name", 10)

	if v.IsValid() {
		t.Error("expected name to be too long")
	}

	user2 := TestUser{Name: "john"}
	v2 := New(&user2)
	v2.MaxLength("Name", 10)

	if !v2.IsValid() {
		t.Error("name should be valid")
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user@domain.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", true}, // empty emails pass (use Required for that)
	}

	for _, tt := range tests {
		user := TestUser{Email: tt.email}
		v := New(&user)
		v.Email("Email")

		if v.IsValid() != tt.valid {
			t.Errorf("email %q: expected valid=%v, got %v", tt.email, tt.valid, v.IsValid())
		}
	}
}

func TestMinInt(t *testing.T) {
	user := TestUser{Age: 15}
	v := New(&user)
	v.MinInt("Age", 18)

	if v.IsValid() {
		t.Error("age should be too low")
	}

	user2 := TestUser{Age: 25}
	v2 := New(&user2)
	v2.MinInt("Age", 18)

	if !v2.IsValid() {
		t.Error("age should be valid")
	}
}

func TestMaxInt(t *testing.T) {
	user := TestUser{Age: 150}
	v := New(&user)
	v.MaxInt("Age", 120)

	if v.IsValid() {
		t.Error("age should be too high")
	}

	user2 := TestUser{Age: 30}
	v2 := New(&user2)
	v2.MaxInt("Age", 120)

	if !v2.IsValid() {
		t.Error("age should be valid")
	}
}

func TestInEnum(t *testing.T) {
	validRoles := []string{"user", "admin", "moderator"}

	user := TestUser{Role: "superuser"}
	v := New(&user)
	v.InEnum("Role", validRoles)

	if v.IsValid() {
		t.Error("role should be invalid")
	}

	user2 := TestUser{Role: "admin"}
	v2 := New(&user2)
	v2.InEnum("Role", validRoles)

	if !v2.IsValid() {
		t.Error("role should be valid")
	}
}

func TestUnique(t *testing.T) {
	checkUnique := func(value interface{}) bool {
		name := value.(string)
		taken := []string{"admin", "root"}
		for _, n := range taken {
			if name == n {
				return false
			}
		}
		return true
	}

	user := TestUser{Name: "admin"}
	v := New(&user)
	v.Unique("Name", checkUnique)

	if v.IsValid() {
		t.Error("name should already be taken")
	}

	user2 := TestUser{Name: "john"}
	v2 := New(&user2)
	v2.Unique("Name", checkUnique)

	if !v2.IsValid() {
		t.Error("name should be unique")
	}
}

func TestCustomMessages(t *testing.T) {
	user := TestUser{Name: ""}
	v := New(&user)
	customMsg := "please enter your name"
	v.Required("Name", customMsg)

	if v.Errors["Name"] != customMsg {
		t.Errorf("expected custom message %q, got %q", customMsg, v.Errors["Name"])
	}
}

func TestMultipleErrors(t *testing.T) {
	user := TestUser{
		Name:     "",
		Email:    "bad-email",
		Password: "123",
		Age:      10,
	}

	v := New(&user)
	v.Required("Name")
	v.Email("Email")
	v.MinLength("Password", 8)
	v.MinInt("Age", 18)

	if v.IsValid() {
		t.Error("expected multiple validation errors")
	}

	if len(v.Errors) != 4 {
		t.Errorf("expected 4 errors, got %d", len(v.Errors))
	}
}

func TestGetErrorsMessages(t *testing.T) {
	user := TestUser{Name: ""}
	v := New(&user)
	v.Required("Name")

	msg := v.GetErrorsMessages()
	if msg == "" {
		t.Error("expected error message")
	}
}

func TestAddError(t *testing.T) {
	user := TestUser{}
	v := New(&user)

	v.AddError("Name", "custom error")

	if v.Errors["Name"] != "custom error" {
		t.Error("expected custom error to be added")
	}

	v.AddError("Name", "second error")
	if v.Errors["Name"] != "custom error" {
		t.Error("expected first error to remain")
	}
}

func TestCheck(t *testing.T) {
	user := TestUser{}
	v := New(&user)

	v.Check(false, "Name", "name is invalid")

	if v.IsValid() {
		t.Error("expected validation to fail")
	}

	if v.Errors["Name"] != "name is invalid" {
		t.Error("expected error message to be set")
	}
}

func TestUniqueSkipsWhenFieldHasError(t *testing.T) {
	user := TestUser{Name: ""}
	v := New(&user)

	v.Required("Name")

	called := false
	checkFunc := func(value interface{}) bool {
		called = true
		return false
	}

	v.Unique("Name", checkFunc)

	if called {
		t.Error("unique check should be skipped when field already has error")
	}
}

func TestFieldNotFound(t *testing.T) {
	user := TestUser{Name: "john"}
	v := New(&user)

	v.Required("NonExistentField")
	v.MinLength("MissingField", 5)
	v.MaxLength("MissingField2", 10)
	v.Email("MissingEmail")
	v.MinInt("MissingAge", 18)
	v.MaxInt("MissingAge2", 100)
	v.InEnum("MissingRole", []string{"admin"})
	v.Unique("MissingName", func(v interface{}) bool { return true })

	fields := []string{
		"NonExistentField", "MissingField", "MissingField2",
		"MissingEmail", "MissingAge", "MissingAge2",
		"MissingRole", "MissingName",
	}
	for _, field := range fields {
		if _, exists := v.Errors[field]; !exists {
			t.Errorf("expected error for missing field %q", field)
		}
	}
}

func TestInvalidType(t *testing.T) {
	user := TestUser{Name: "john", Age: 25}
	v := New(&user)

	// MinInt/MaxInt on a string field
	v.MinInt("Name", 5)
	v.MaxInt("Name", 50)

	// MinLength/MaxLength on an int field
	v.MinLength("Age", 2)
	v.MaxLength("Age", 10)

	// Email on an int field
	v.Email("Age")

	// InEnum on an int field
	v.InEnum("Age", []string{"admin"})

	invalidFields := []string{"Name", "Age"}
	for _, field := range invalidFields {
		if _, exists := v.Errors[field]; !exists {
			t.Errorf("expected invalid-type error for field %q", field)
		}
	}
}
