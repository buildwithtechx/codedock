package transports

import (
	"errors"
	"strings"
)

// VerifyTeamsRequest validates an incoming Microsoft Teams / Bot Framework request.
func VerifyTeamsRequest(configuredAppID, authorizationHeader string) error {
	if configuredAppID == "" {
		return nil
	}
	if authorizationHeader == "" {
		return errors.New("missing authorization header from teams request")
	}
	if !strings.HasPrefix(strings.ToLower(authorizationHeader), "bearer ") {
		return errors.New("invalid authorization header scheme, expected bearer token")
	}
	return nil
}
