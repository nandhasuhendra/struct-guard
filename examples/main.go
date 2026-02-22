package main

import (
	"fmt"

	validator "github.com/nandhasuhendra/struct-guard"
)

type User struct {
	Username string
	Email    string
	Password string
	Age      int
	Role     string
}

func isUsernameUnique(value interface{}) bool {
	username := value.(string)
	existingUsers := []string{"admin", "test", "user123"}
	for _, u := range existingUsers {
		if username == u {
			return false
		}
	}
	return true
}

// validUser demonstrates a fully valid submission that passes all rules.
func validUser() {
	fmt.Println("=== Valid user ===")

	user := User{
		Username: "john",
		Email:    "john@example.com",
		Password: "mypassword123",
		Age:      25,
		Role:     "user",
	}

	v := validator.New(&user)
	v.Required("Username")
	v.Required("Email")
	v.Required("Password")
	v.MinLength("Username", 3)
	v.MaxLength("Username", 20)
	v.Email("Email")
	v.MinLength("Password", 8)
	v.MinInt("Age", 18)
	v.MaxInt("Age", 120)
	v.InEnum("Role", []string{"user", "admin", "moderator"})
	v.Unique("Username", isUsernameUnique)

	if !v.IsValid() {
		fmt.Println("Validation errors:")
		fmt.Print(v.GetErrorsMessages())
		return
	}
	fmt.Println("User validated successfully")
}

// invalidUser demonstrates common validation failures and the clear error
// messages that are returned to the caller instead of panicking.
func invalidUser() {
	fmt.Println("\n=== Invalid user ===")

	user := User{
		Username: "",            // fails Required
		Email:    "not-an-email", // fails Email format
		Password: "short",       // fails MinLength(8)
		Age:      15,             // fails MinInt(18)
		Role:     "superuser",   // fails InEnum
	}

	v := validator.New(&user)
	v.Required("Username")
	v.Email("Email")
	v.MinLength("Password", 8)
	v.MinInt("Age", 18)
	v.InEnum("Role", []string{"user", "admin", "moderator"})

	fmt.Println("Validation errors:")
	fmt.Print(v.GetErrorsMessages())
}

// fieldErrors demonstrates that unknown field names and wrong types produce
// descriptive errors instead of panicking.
func fieldErrors() {
	fmt.Println("\n=== Field / type errors ===")

	user := User{Username: "alice", Age: 30}
	v := validator.New(&user)

	// Field does not exist on the struct — clear error, no panic.
	v.Required("PhoneNumber")

	// MinInt called on a string field — type mismatch error, no panic.
	v.MinInt("Username", 1)

	// Email called on an int field — type mismatch error, no panic.
	v.Email("Age")

	fmt.Println("Errors caught without panicking:")
	fmt.Print(v.GetErrorsMessages())
}

func main() {
	validUser()
	invalidUser()
	fieldErrors()
}
