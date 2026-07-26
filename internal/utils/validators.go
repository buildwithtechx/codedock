package utils

import (
	"errors"
	"net/url"
	"strings"
)

func ValidateURL(inputURL string) (string, error) {
	trimmedURL := strings.TrimSpace(inputURL)
	if trimmedURL == "" {
		return "", errors.New("url cannot be empty or whitespace")
	}
	u, err := url.ParseRequestURI(trimmedURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", errors.New("invalid url: must be an absolute http(s) url")
	}
	return trimmedURL, nil
}

func ValidateEventTypes(eventTypes []string) ([]string, error) {
	if eventTypes == nil {
		return nil, nil
	}
	validated := make([]string, len(eventTypes))
	for i, et := range eventTypes {
		trimmed := strings.TrimSpace(et)
		if trimmed == "" {
			return nil, errors.New("event type cannot be empty or whitespace")
		}
		if strings.Contains(trimmed, ",") {
			return nil, errors.New("event type cannot contain commas")
		}
		validated[i] = trimmed
	}
	return validated, nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
