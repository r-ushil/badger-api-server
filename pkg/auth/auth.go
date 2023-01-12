package auth

import (
	"badger-api/pkg/server"
	"context"
	"errors"
	"strings"
)

func ParseIdToken(s *server.ServerContext, idToken string) (string, error) {
	auth := s.GetFirebaseAuth()

	decodedToken, err := auth.VerifyIDToken(context.Background(), idToken)

	if err != nil {
		return "", err
	}

	return decodedToken.UID, nil
}

func ExtractBearerToken(authHeader string) (string, error) {
	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		return "", errors.New("Missing bearer token")
	}

	idToken := splitToken[1]

	return idToken, nil
}

func ParseAuthHeader(s *server.ServerContext, authHeader string) (string, error) {
	idToken, err := ExtractBearerToken(authHeader)

	if err == nil {
		return ParseIdToken(s, idToken)
	}

	// If not bearer token, then authorization is the UserID
	userId := authHeader

	return userId, nil
}
