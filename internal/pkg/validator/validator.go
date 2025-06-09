package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator membungkus instance validator
type Validator struct {
	validator *validator.Validate
}

// ValidationError merepresentasikan error validasi tunggal
type ValidationError struct {
	Field   string `json:"field"`   // Nama field yang error
	Tag     string `json:"tag"`     // Tag validasi yang gagal
	Value   string `json:"value"`   // Nilai yang gagal divalidasi
	Message string `json:"message"` // Pesan error yang mudah dibaca
}

// ValidationErrors koleksi error validasi
type ValidationErrors []ValidationError

// Error mengubah error validasi menjadi string
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// New membuat instance validator baru dengan konfigurasi default
func New() *Validator {
	v := validator.New()

	// Registrasi validasi kustom
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("username", validateUsername)

	// Menggunakan nama json tag sebagai nama field
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validator: v}
}

// ValidateStruct memvalidasi struktur data
func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			var validationErrors ValidationErrors
			for _, err := range valErrs {
				validationErrors = append(validationErrors, ValidationError{
					Field:   err.Field(),
					Tag:     err.Tag(),
					Value:   fmt.Sprintf("%v", err.Value()),
					Message: getErrorMessage(err),
				})
			}
			return validationErrors
		}
		return err
	}
	return nil
}

// validatePassword memvalidasi kekuatan password
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// validatePhone memvalidasi format nomor telepon
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(digits) < 10 || len(digits) > 15 {
		return false
	}

	// Format Indonesia
	if strings.HasPrefix(digits, "62") {
		return len(digits) >= 11 && len(digits) <= 13
	}

	// Format internasional
	if strings.HasPrefix(digits, "+") {
		return len(digits) >= 11 && len(digits) <= 16
	}

	// Format lokal
	if strings.HasPrefix(digits, "0") {
		return len(digits) >= 10 && len(digits) <= 12
	}

	return true
}

// validateUsername memvalidasi format username
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 30 {
		return false
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return false
	}

	return !(strings.HasPrefix(username, "_") || strings.HasPrefix(username, "-") ||
		strings.HasSuffix(username, "_") || strings.HasSuffix(username, "-"))
}

// getErrorMessage menghasilkan pesan error yang mudah dibaca
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "email":
		return fmt.Sprintf("%s harus berupa alamat email yang valid", field)
	case "min":
		return fmt.Sprintf("%s minimal harus %s karakter", field, err.Param())
	case "max":
		return fmt.Sprintf("%s maksimal harus %s karakter", field, err.Param())
	case "password":
		return fmt.Sprintf("%s harus minimal 8 karakter dengan huruf besar, kecil, angka, dan karakter khusus", field)
	case "phone":
		return fmt.Sprintf("%s harus berupa nomor telepon yang valid", field)
	case "username":
		return fmt.Sprintf("%s harus 3-30 karakter, alfanumerik dengan underscore/hyphen (tidak di awal/akhir)", field)
	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}

// ToMap mengubah error validasi menjadi map[field][]message
func (ve ValidationErrors) ToMap() map[string][]string {
	result := make(map[string][]string)
	for _, err := range ve {
		result[err.Field] = append(result[err.Field], err.Message)
	}
	return result
}

// IsValidationError mengecek apakah error merupakan error validasi
func IsValidationError(err error) bool {
	_, ok := err.(ValidationErrors)
	return ok
}