# struct-guard

**struct-guard** is a lightweight, dependency-free validation library for Go. 

It brings the expressive, declarative style of Rails Active Record validations to Go structs, without the magic of callbacks or tight coupling to a database. It is designed specifically for **Domain-Driven Design (DDD)**, allowing you to validate business rules (like uniqueness) without polluting your domain entities with database logic.

## Why struct-guard?

Go validation libraries often fall into two traps:
1. **Tag Soup:** Struct tags (`binding:"required"`) clutter your domain models and couple them to specific validation logic.
2. **Boilerplate:** Writing `if len(str) == 0` fifty times is tedious and prone to error.

**struct-guard** offers a middle ground:
* **Fluent API:** Readable, chainable rules.
* **Dependency Injection:** Handles stateful checks (like "is email unique?") via closures, keeping your structs pure.
* **Zero Dependencies:** Uses only the Go standard library (`reflect`, `net/mail`).

## Installation

```bash
go get github.com/nandhasuhendra/struct-guard
```

## Usage

```go
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
    v.MinLength("Username", 3)
    v.MaxLength("Username", 20)
    v.Email("Email")
    v.MinLength("Password", 8)
    v.MinInt("Age", 18)
    v.InEnum("Role", []string{"user", "admin", "moderator"})

    if !v.IsValid() {
        fmt.Println("Validation errors:")
        fmt.Print(v.GetErrorsMessages())
        return
    }

    fmt.Println("User validated successfully")
}
```

## Available Validators

### Required
Checks if a field is not zero-value (empty string, 0, nil, etc.)
```go
v.Required("Username")
v.Required("Email", "custom error message")
```

### MinLength / MaxLength
Validates string length
```go
v.MinLength("Password", 8)
v.MaxLength("Username", 20)
```

### MinInt / MaxInt
Validates integer range
```go
v.MinInt("Age", 18)
v.MaxInt("Age", 120)
```

### Email
Validates email format using `net/mail`
```go
v.Email("Email")
```

### InEnum
Validates against a list of allowed values
```go
v.InEnum("Role", []string{"user", "admin", "moderator"})
```

### Unique
Custom validation with a callback function (useful for database checks)
```go
checkUnique := func(value interface{}) bool {
    username := value.(string)
    // Check database or cache
    return !userExists(username)
}

v.Unique("Username", checkUnique)
```

### Custom Validation
Use `Check` for custom validation logic
```go
v.Check(user.Age >= 18, "Age", "must be 18 or older")
```

## Custom Error Messages

All validators accept an optional custom error message:
```go
v.Required("Username", "Please provide a username")
v.MinLength("Password", 8, "Password too short")
v.Email("Email", "That doesn't look like a valid email")
```

## Error Handling

```go
v := validator.New(&user)
v.Required("Username")
v.Email("Email")

// Check if validation passed
if !v.IsValid() {
    // Get all error messages as string
    fmt.Print(v.GetErrorsMessages())
    
    // Or access individual errors
    for field, msg := range v.Errors {
        fmt.Printf("%s: %s\n", field, msg)
    }
}
```

## License

MIT License - see [LICENSE](LICENSE) file for details
