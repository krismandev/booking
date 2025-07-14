package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const tagCustom = "errormgs"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func isValidCardNumber(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	re := regexp.MustCompile(`^\d{16}$`)
	return re.MatchString(cardNumber)
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// At least 1 uppercase letter
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// At least 1 lowercase letter
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	// At least 1 special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	return hasUppercase && hasLowercase && hasSpecial
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("cardnumber", isValidCardNumber)
	v.RegisterValidation("password", passwordValidator)

	return &CustomValidator{Validator: v}
}

func errorTagFunc[T interface{}](obj interface{}, snp string, fieldname, actualTag string) error {
	o := obj.(T)

	if !strings.Contains(snp, fieldname) {
		return nil
	}

	fieldArr := strings.Split(snp, ".")
	rsf := reflect.TypeOf(o)

	for i := 1; i < len(fieldArr); i++ {
		field, found := rsf.FieldByName(fieldArr[i])
		if found {
			if fieldArr[i] == fieldname {
				customMessage := field.Tag.Get(tagCustom)
				if customMessage != "" {
					return fmt.Errorf("%s: %s (%s)", fieldname, customMessage, actualTag)
				}
				return nil
			} else {
				if field.Type.Kind() == reflect.Ptr {
					// If the field type is a pointer, dereference it
					rsf = field.Type.Elem()
				} else {
					rsf = field.Type
				}
			}
		}
	}
	return nil
}

func ValidateFunc[T interface{}](obj interface{}, validate *validator.Validate) (errs error) {
	o := obj.(T)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Validate:", r)
			errs = fmt.Errorf("can't validate %+v", r)
		}
	}()

	if err := validate.Struct(o); err != nil {
		errorValid := err.(validator.ValidationErrors)
		for _, e := range errorValid {
			// snp  X.Y.Z
			snp := e.StructNamespace()
			errmgs := errorTagFunc[T](obj, snp, e.Field(), e.ActualTag())
			if errmgs != nil {
				errs = errors.Join(errs, fmt.Errorf("%w", errmgs))
			} else {
				errs = errors.Join(errs, fmt.Errorf("%w", e))
			}
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

func FormatValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldErr := range validationErrors {
			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s field is %s", fieldErr.StructNamespace(), fieldErr.Tag()))
			case "number":
				messages = append(messages, fmt.Sprintf("%s should be number", fieldErr.Field()))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s character / item ", fieldErr.Field(), fieldErr.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s cannot exceed %s characters", fieldErr.Field(), fieldErr.Param()))
			case "password":
				messages = append(messages, fmt.Sprintf("%s must contain at least one uppercase letter, one lowercase letter, and one special character", fieldErr.Field()))
			default:
				messages = append(messages, fmt.Sprintf("Wrong %s format", fieldErr.StructNamespace()))
			}
		}
		return strings.Join(messages, ", ")
	}
	return "Invalid request"
}
