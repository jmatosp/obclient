package authorization

import (
	"github.com/pkg/errors"
)

type Authenticator interface {
	Authenticate() (Token, error)
}

type authenticator struct {
	credentialsGranter CredentialsGranter
	accessConsenter    AccessConsenter
	psuAccessConsenter PSUAccessConsenter
	tokenGenerator     TokenGenerator
}

func NewAuthenticator(
	credentialsGranter CredentialsGranter,
	accessConsenter AccessConsenter,
	psuAccessConsenter PSUAccessConsenter,
	generator TokenGenerator,
) Authenticator {
	return authenticator{
		credentialsGranter: credentialsGranter,
		accessConsenter:    accessConsenter,
		psuAccessConsenter: psuAccessConsenter,
		tokenGenerator:     generator,
	}
}

func (a authenticator) Authenticate() (Token, error) {
	grantsToken, err := a.credentialsGranter.Request()
	if err != nil {
		return NoToken, errors.Wrap(err, "error authenticating")
	}

	accessConsent, err := a.accessConsenter.Request(grantsToken)
	if err != nil {
		return NoToken, errors.Wrap(err, "error authenticating")
	}

	code, err := a.psuAccessConsenter.Request(accessConsent)
	if err != nil {
		return NoToken, errors.Wrap(err, "error authenticating")
	}

	token, err := a.tokenGenerator.Request(code)
	if err != nil {
		return NoToken, errors.Wrap(err, "error authenticating")
	}

	return token, nil
}
