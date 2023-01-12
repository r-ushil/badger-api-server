package auth

import (
	"badger-api/pkg/server"
	"context"
)

func ParseIdToken(s *server.ServerContext, idToken string) (string, error) {
	auth := s.GetFirebaseAuth()

	decodedToken, err := auth.VerifyIDToken(context.Background(), idToken)

	if err != nil {
		return "", err
	}

	return decodedToken.Subject, nil
}
