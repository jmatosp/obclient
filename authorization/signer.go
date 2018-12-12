package authorization

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Signer interface {
	Sign(claims jwt.Claims) (string, error)
}

func NewRS256Singer(certificates Certificate) Signer {
	return &rs256Signer{
		certs: certificates,
	}
}

type rs256Signer struct {
	certs Certificate
}

func (s *rs256Signer) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := s.certs.PrivateKey()
	if err != nil {
		return "", errors.Wrap(err, "error signing claims")
	}

	return token.SignedString(privateKey)
}
