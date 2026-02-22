package validator

import (
	"fmt"
	"net/mail"
	"reflect"
	"strings"
)

const (
	REQUIRED_FIELD  = "required"
	MIN_LENGTH      = "too short, minimum %d characters"
	MAX_LENGTH      = "too long, maximum %d characters"
	MIN_VALUE       = "too small, minimum value is %d"
	MAX_VALUE       = "too large, maximum value is %d"
	EMAIL_FORMAT    = "invalid email address"
	ENUM_OPTION     = "must be one of: %s"
	UNIQUE_CHECK    = "already taken"
	EQUAL_VALUE     = "must be equal to %v"
	NOT_EQUAL_VALUE = "must not be equal to %v"

	FIELD_NOT_FOUND = "unknown field '%s'"
	INVALID_TYPE    = "unsupported type for field '%s'"

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

func (v *Validator) Check(equal bool, field, message string) {
	if !equal {
		v.AddError(field, message)
	}
}

// -- Validation helper functions --

// Required checks if the field is not zero-value (empty string, 0, nil, etc.)
func (v *Validator) Required(field string, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	msg := msgOr(REQUIRED_FIELD, customMsg...)
	v.Check(!val.IsZero(), field, msg)
}

func (v *Validator) MinInt(field string, min int, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// valid integer kinds
	default:
		v.AddError(field, fmt.Sprintf(INVALID_TYPE, field))
		return
	}
	msg := msgOr(fmt.Sprintf(MIN_VALUE, min), customMsg...)
	v.Check(val.Int() >= int64(min), field, msg)
}

func (v *Validator) MaxInt(field string, max int, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// valid integer kinds
	default:
		v.AddError(field, fmt.Sprintf(INVALID_TYPE, field))
		return
	}
	msg := msgOr(fmt.Sprintf(MAX_VALUE, max), customMsg...)
	v.Check(val.Int() <= int64(max), field, msg)
}

func (v *Validator) MinLength(field string, min int, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	switch val.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		// valid kinds that support Len()
	default:
		v.AddError(field, fmt.Sprintf(INVALID_TYPE, field))
		return
	}
	msg := msgOr(fmt.Sprintf(MIN_LENGTH, min), customMsg...)
	v.Check(val.Len() >= min, field, msg)
}

func (v *Validator) MaxLength(field string, max int, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	switch val.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		// valid kinds that support Len()
	default:
		v.AddError(field, fmt.Sprintf(INVALID_TYPE, field))
		return
	}
	msg := msgOr(fmt.Sprintf(MAX_LENGTH, max), customMsg...)
	v.Check(val.Len() <= max, field, msg)
}

func (v *Validator) Email(field string, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	if val.Kind() != reflect.String {
		v.AddError(field, fmt.Sprintf(INVALID_TYPE, field))
		return
	}
	msg := msgOr(EMAIL_FORMAT, customMsg...)
	if val.String() == "" {
		return
	}
	_, err := mail.ParseAddress(val.String())
	v.Check(err == nil, field, msg)
}

func (v *Validator) InEnum(field string, validValues []string, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	if val.Kind() != reflect.String {
		v.AddError(field, fmt.Sprintf(INVALID_TYPE, field))
		return
	}
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

// Equal checks if the field value equals the given value.
// It supports any comparable type (string, int, float, bool, etc.).
func (v *Validator) Equal(field string, expected interface{}, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	defaultMsg := fmt.Sprintf(EQUAL_VALUE, expected)
	msg := msgOr(defaultMsg, customMsg...)
	v.Check(reflect.DeepEqual(val.Interface(), expected), field, msg)
}

// NotEqual checks if the field value does not equal the given value.
// It supports any comparable type (string, int, float, bool, etc.).
func (v *Validator) NotEqual(field string, unexpected interface{}, customMsg ...string) {
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	defaultMsg := fmt.Sprintf(NOT_EQUAL_VALUE, unexpected)
	msg := msgOr(defaultMsg, customMsg...)
	v.Check(!reflect.DeepEqual(val.Interface(), unexpected), field, msg)
}

func (v *Validator) Unique(field string, checkFunc func(interface{}) bool, customMsg ...string) {
	if _, hasError := v.Errors[field]; hasError {
		return
	}
	val, ok := v.val(field)
	if !ok {
		v.AddError(field, fmt.Sprintf(FIELD_NOT_FOUND, field))
		return
	}
	isUnique := checkFunc(val.Interface())
	msg := msgOr(UNIQUE_CHECK, customMsg...)
	v.Check(isUnique, field, msg)
}

// -- Private Helper functions --

// val reflects the value from the struct by field name.
// Returns the reflected value and false (instead of panicking) when the field does not exist.
func (v *Validator) val(field string) (reflect.Value, bool) {
	f := v.data.FieldByName(field)
	if !f.IsValid() {
		return reflect.Value{}, false
	}
	return f, true
}

func msgOr(defaultMsg string, custom ...string) string {
	if len(custom) > 0 && custom[0] != "" {
		return custom[0]
	}
	return defaultMsg
}
