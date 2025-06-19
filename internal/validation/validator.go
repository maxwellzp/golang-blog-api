package validation

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{
		validate: validator.New(),
	}

	// Register custom validations
	v.validate.RegisterValidation("containsuppercase", containsUpperCase)
	v.validate.RegisterValidation("containslowercase", containsLowerCase)
	v.validate.RegisterValidation("containsnumber", containsNumber)
	v.validate.RegisterValidation("containsspecial", containsSpecial)

	return v
}

func (v *Validator) ValidateStruct(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := toJSONFieldName(err.StructField(), s)
		errors[field] = validationMessage(err)
	}

	return errors
}

func toJSONFieldName(fieldName string, s interface{}) string {
	rt := reflect.TypeOf(s)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	field, found := rt.FieldByName(fieldName)
	if !found {
		return strings.ToLower(fieldName)
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return strings.ToLower(fieldName)
	}
	return strings.Split(jsonTag, ",")[0]
}

func validationMessage(fe validator.FieldError) string {
	var msg string
	switch fe.Tag() {
	case "required":
		msg = fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return "must be a valid email"
	case "min":
		msg = fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		msg = fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "containsuppercase":
		msg = fmt.Sprintf("%s must contain at least one uppercase letter", fe.Field())
	case "containslowercase":
		msg = fmt.Sprintf("%s must contain at least one lowercase letter", fe.Field())
	case "containsnumber":
		msg = fmt.Sprintf("%s must contain at least one number", fe.Field())
	case "containsspecial":
		msg = fmt.Sprintf("%s must contain at least one special character", fe.Field()) + ` (!@#$%^&*)`
	default:
		msg = fmt.Sprintf("%s is invalid", fe.Field())
	}
	return msg
}

func containsUpperCase(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func containsLowerCase(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func containsNumber(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func containsSpecial(fl validator.FieldLevel) bool {
	specialChars := "!@#$%^&*"
	fieldValue := fl.Field().String()

	for _, char := range specialChars {
		if strings.ContainsRune(fieldValue, char) {
			return true
		}
	}
	return false
}
