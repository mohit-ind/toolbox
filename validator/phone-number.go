package validator

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// ParseDutchPhoneNumber parses the supplied phoneNumber string by
// removing all spaces, hypens and the +31 prefix, and adds a 0 instead.
// It will return the parsed string, and an optional error,
// if the parsed string is not a Dutch Mobile Phone Number.
// e.g.: +31 6-3448-2527 -> 0634482527, nil
func ParseDutchPhoneNumber(phoneNumber string) (string, error) {
	// Remove spaces
	parsed := strings.ReplaceAll(phoneNumber, " ", "")
	// Remove hyphens
	parsed = strings.ReplaceAll(parsed, "-", "")

	// If the number is prefixed with +31
	if strings.HasPrefix(parsed, "+31") {
		// Remove the prefix
		parsed = strings.TrimPrefix(parsed, "+31")
		// Assert it to pattern: 612345678
		phoneRegex := regexp.MustCompile(`^6[0-9]{8}$`)
		if !phoneRegex.MatchString(parsed) {
			return "", errors.Errorf("%s is not a Dutch Mobile Number", phoneNumber)
		}
		// Prefix with 0
		parsed = "0" + parsed
	}

	// Now it must match 0612345678
	phoneRegex := regexp.MustCompile(`^06[0-9]{8}$`)
	if !phoneRegex.MatchString(parsed) {
		return "", errors.Errorf("%s is not a Dutch Mobile Number", phoneNumber)
	}

	return parsed, nil
}
