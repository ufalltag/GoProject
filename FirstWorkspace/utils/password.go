package utils

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

// ValidatePasswordStrength проверяет сложность пароля.
// Пароль должен быть не короче 8 символов и содержать
// заглавную и строчную буквы, цифру и спецсимвол.
func ValidatePasswordStrength(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return errors.New("Пароль должен содержать минимум 8 символов")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("Пароль должен содержать заглавную букву")
	}
	if !hasLower {
		return errors.New("Пароль должен содержать строчную букву")
	}
	if !hasDigit {
		return errors.New("Пароль должен содержать цифру")
	}
	if !hasSpecial {
		return errors.New("Пароль должен содержать спецсимвол")
	}

	return nil
}
