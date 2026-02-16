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

func main() {
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
	v.InEnum("Role", []string{"user", "admin", "moderator"})
	v.Unique("Username", isUsernameUnique)

	if !v.IsValid() {
		fmt.Println("Validation errors:")
		fmt.Print(v.GetErrorsMessages())
		return
	}

	fmt.Println("User validated successfully")
}
