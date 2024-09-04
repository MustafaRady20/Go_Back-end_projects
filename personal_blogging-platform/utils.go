package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"unicode"
)

const emailEgexp = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(emailEgexp)
	return emailRegex.MatchString(e)
}

func isPasswordValid(p string) bool {
	var hasUpperLetter, hasSpecialChar, hasNumber, length bool
	for _, ch := range p {
		if unicode.IsUpper(ch) {
			hasUpperLetter = true
		}
		if unicode.IsDigit(ch) {
			hasNumber = true
		}
		if unicode.IsSymbol(ch) {
			hasSpecialChar = true
		}

	}
	if len(p) >= 8 {
		length = true
	}

	if hasNumber && hasSpecialChar && hasUpperLetter && length {
		return true
	}
	return true
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
