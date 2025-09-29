package pstring

import (
	"errors"
	"fmt"
	"regexp"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	b := re.Find([]byte(email))
	if b == nil {
		msg := fmt.Sprintf("\"%s\" is not a valid email address", email)
		err := errors.New(msg)
		return err
	}
	return nil
}
