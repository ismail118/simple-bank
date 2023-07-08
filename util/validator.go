package util

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-z0-9\s]+$`).MatchString
)

func ValidateString(s string, minLength, maxLength int) error {
	n := len(s)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(s string) error {
	err := ValidateString(s, 3, 100)
	if err != nil {
		return err
	}

	if !isValidUsername(s) {
		return fmt.Errorf("must contain only lowercase latters, digits or underscore")
	}

	return nil
}

func ValidatePassword(s string) error {
	return ValidateString(s, 6, 100)
}

func ValidateEmail(s string) error {
	err := ValidateString(s, 3, 100)
	if err != nil {
		return err
	}

	_, err = mail.ParseAddress(s)
	if err != nil {
		return fmt.Errorf("is not valid email addreess")
	}

	return nil
}

func ValidateFullName(s string) error {
	err := ValidateString(s, 3, 100)
	if err != nil {
		return err
	}

	if !isValidFullName(s) {
		return fmt.Errorf("must contain only latters, digits or underscore")
	}

	return nil
}
