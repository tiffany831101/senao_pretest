package utils

import (
	"unicode"
	"unicode/utf8"
)

func IsPasswordComplex(password string) bool {

	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	if utf8.RuneCountInString(password) > 32 {
		return false
	}

	var hasUpperCase, hasLowerCase, hasNumber bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasUpperCase && hasLowerCase && hasNumber
}
