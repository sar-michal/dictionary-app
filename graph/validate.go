package graph

import (
	"fmt"
	"regexp"
	"strings"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

// sanitizeInput trims leading/trailing whitespace and collapses multiple spaces.
func sanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return whitespaceRegex.ReplaceAllString(trimmed, " ")
}

// validateNonEmpty ensures that the input is not empty.
func validateNonEmpty(fieldValue string) error {
	if strings.TrimSpace(fieldValue) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

// validateLength checks if the input length is within the specified maximum.
func validateLength(fieldValue string, max int) error {
	length := len(fieldValue)
	if length > max {
		return fmt.Errorf("input must be at most %d characters", max)
	}
	return nil
}

// validateInput sanitizes the input and checks that it's non-empty and within the maximum length.
// It returns the sanitized string or an error if validation fails.
func validateInput(input string) (string, error) {
	sanitized := sanitizeInput(input)
	const maxLength int = 200
	if err := validateNonEmpty(sanitized); err != nil {
		return "", err
	}

	if err := validateLength(sanitized, maxLength); err != nil {
		return "", err
	}

	return sanitized, nil
}
