package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[\p{L}\p{N}._%+-]+@[\p{L}\p{N}.-]+\.[\p{L}]{2,}$`)
	return re.MatchString(email)
}

func validatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func validateLoginRequest(req *LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if !validateEmail(req.Email) {
		return errors.New("невалидный формат email")
	}
	if !validatePassword(req.Password) {
		return errors.New("невалидные учетные данные")
	}
	return nil
}

func validateRegisterRequest(req *RegisterRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if !validateEmail(req.Email) {
		return errors.New("невалидный формат email")
	}
	if !validatePassword(req.Password) {
		return errors.New("пароль должен быть не менее 6 символов, иметь разный регистр, иметь хотя бы одну цифру")
	}
	return nil
}

func parseIntQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}