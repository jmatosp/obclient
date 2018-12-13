package authorization

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"strings"
)

type Signer interface {
	Sign(claims jwt.Claims) (string, error)
}

// NewSigner returns a Signer object given a possibility of signing algo
// if no compatible sign algo is found an error is returned
func NewSigner(certificate Certificate, methods []string) (Signer, error) {
	for _, method := range methods {
		switch strings.ToLower(method) {
		case "rs256":
			return NewSingerWithMethod(certificate, jwt.SigningMethodRS256), nil
		case "ps256":
			return NewSingerWithMethod(certificate, jwt.SigningMethodPS256), nil
		case "hs256":
			return NewSingerWithMethod(certificate, jwt.SigningMethodHS256), nil
		case "rs384":
			return NewSingerWithMethod(certificate, jwt.SigningMethodRS384), nil
		case "ps384":
			return NewSingerWithMethod(certificate, jwt.SigningMethodPS384), nil
		case "hs384":
			return NewSingerWithMethod(certificate, jwt.SigningMethodHS384), nil
		case "rs512":
			return NewSingerWithMethod(certificate, jwt.SigningMethodRS512), nil
		case "ps512":
			return NewSingerWithMethod(certificate, jwt.SigningMethodPS512), nil
		case "hs512":
			return NewSingerWithMethod(certificate, jwt.SigningMethodHS512), nil
		}
	}
	return nil, errors.New("error could not find a compatible signing method")
}

func NewSingerWithMethod(certificates Certificate, method jwt.SigningMethod) Signer {
	return &signer{
		certs:      certificates,
		signMethod: method,
	}
}

type signer struct {
	certs      Certificate
	signMethod jwt.SigningMethod
}

func (s *signer) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(s.signMethod, claims)

	privateKey, err := s.certs.PrivateKey()
	if err != nil {
		return "", errors.Wrap(err, "error signing claims")
	}

	return token.SignedString(privateKey)
}
