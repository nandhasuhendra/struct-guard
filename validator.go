package validator

import (
	"fmt"
	"net/mail"
	"reflect"
	"strings"
)

const (
	REQUIRED_FIELD = "this field is required"
	MIN_LENGTH     = "must be at least %d characters"
	MAX_LENGTH     = "must be at most %d characters"
	EMAIL_FORMAT   = "invalid email format"
	ENUM_OPTION    = "must be one of the allowed options: %s"
	UNIQUE_CHECK   = "must be unique" // Added constant for uniqueness

	DEFAULT_MIN_LENGTH = 8
	DEFAULT_MAX_LENGTH = 72
)

type ValidationErrors map[string]string

type Validator struct {
	data   reflect.Value
	Errors ValidationErrors
}

func New(s interface{}) *Validator {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return &Validator{
		data:   val,
		Errors: make(ValidationErrors),
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) GetErrorsMessages() string {
	messages := ""
	for field, msg := range v.Errors {
		messages += fmt.Sprintf("%s: %s\n", field, msg)
	}
	return messages
}

func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// -- Validation helper functions --

// Required checks if the field is not zero-value (empty string, 0, nil, etc.)
func (v *Validator) Required(field string, customMsg ...string) {
	val := v.val(field)
	msg := msgOr(REQUIRED_FIELD, customMsg...)
	v.Check(!val.IsZero(), field, msg)
}

func (v *Validator) MinInt(field string, min int, customMsg ...string) {
	val := v.val(field)
	msg := msgOr(fmt.Sprintf(MIN_LENGTH, min), customMsg...)

	// Convert to int64 for safe comparison
	currentVal := val.Int()
	v.Check(currentVal >= int64(min), field, msg)
}

func (v *Validator) MaxInt(field string, max int, customMsg ...string) {
	val := v.val(field)
	msg := msgOr(fmt.Sprintf(MAX_LENGTH, max), customMsg...)

	currentVal := val.Int()
	v.Check(currentVal <= int64(max), field, msg)
}

func (v *Validator) MinLength(field string, min int, customMsg ...string) {
	val := v.val(field)
	msg := msgOr(fmt.Sprintf(MIN_LENGTH, min), customMsg...)
	v.Check(val.Len() >= min, field, msg)
}

func (v *Validator) MaxLength(field string, max int, customMsg ...string) {
	val := v.val(field)
	msg := msgOr(fmt.Sprintf(MAX_LENGTH, max), customMsg...)
	v.Check(val.Len() <= max, field, msg)
}

func (v *Validator) Email(field string, customMsg ...string) {
	val := v.val(field)
	msg := msgOr(EMAIL_FORMAT, customMsg...)

	if val.String() == "" {
		return
	}

	_, err := mail.ParseAddress(val.String())
	v.Check(err == nil, field, msg)
}

func (v *Validator) InEnum(field string, validValues []string, customMsg ...string) {
	val := v.val(field)
	currentVal := val.String()

	found := false
	for _, valid := range validValues {
		if currentVal == valid {
			found = true
			break
		}
	}

	defaultMsg := fmt.Sprintf(ENUM_OPTION, strings.Join(validValues, ", "))
	msg := msgOr(defaultMsg, customMsg...)

	v.Check(found, field, msg)
}

func (v *Validator) Unique(field string, checkFunc func(interface{}) bool, customMsg ...string) {
	if _, hasError := v.Errors[field]; hasError {
		return
	}

	val := v.val(field)
	isUnique := checkFunc(val.Interface())

	msg := msgOr(UNIQUE_CHECK, customMsg...)
	v.Check(isUnique, field, msg)
}

// -- Private Helper functions --

// val reflects the value from the struct by field name
func (v *Validator) val(field string) reflect.Value {
	f := v.data.FieldByName(field)
	if !f.IsValid() {
		panic(fmt.Sprintf("validator: field '%s' not found in struct", field))
	}
	return f
}

func msgOr(defaultMsg string, custom ...string) string {
	if len(custom) > 0 && custom[0] != "" {
		return custom[0]
	}
	return defaultMsg
}
