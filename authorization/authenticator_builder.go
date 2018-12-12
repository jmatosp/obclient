package authorization

import (
	"errors"
)

type AuthenticatorBuilder struct {
	client                Client
	fapiFinancialId       string
	accessConsentEndpoint string
	wellKnownEndpoint     string
	redirectUrl           string
	certFile              string
	keyFile               string
	rootCAs               []string
}

func NewAuthenticatorBuilder() *AuthenticatorBuilder {
	return &AuthenticatorBuilder{}
}

func (c *AuthenticatorBuilder) Build() (Authenticator, error) {
	err := c.mustValidate()

	config, err := GetConfiguration(c.wellKnownEndpoint)
	if err != nil {
		return nil, err
	}

	return NewAuthenticator(
		c.makeCredentialsGranter(config),
		c.makeAccessConsenter(),
		c.makePSUAccessConsenter(),
		c.makeTokenGenerator(config),
	), nil
}

func (c *AuthenticatorBuilder) mustValidate() error {
	if c.client.Id == "" || c.client.Secret == "" {
		return errors.New("error client not provided")
	}

	if c.fapiFinancialId == "" {
		return errors.New("error fapiFinancialId not provided")
	}

	if c.accessConsentEndpoint == "" {
		return errors.New("error accessConsentEndpoint not provided")
	}

	if c.wellKnownEndpoint == "" {
		return errors.New("error wellKnownEndpoint not provided")
	}

	if c.redirectUrl == "" {
		return errors.New("error redirectUrl not provided")
	}

	if c.certFile == "" {
		return errors.New("error certFile not provided")
	}

	if c.keyFile == "" {
		return errors.New("error keyFile not provided")
	}

	if len(c.rootCAs) == 0 {
		return errors.New("error need at lease one rootCA")
	}

	return nil
}

func (c *AuthenticatorBuilder) WithClient(client Client) *AuthenticatorBuilder {
	c.client = client
	return c
}

func (c *AuthenticatorBuilder) WithFapiFinancialId(id string) *AuthenticatorBuilder {
	c.fapiFinancialId = id
	return c
}

func (c *AuthenticatorBuilder) WithAccessConsentEndpoint(endpoint string) *AuthenticatorBuilder {
	c.accessConsentEndpoint = endpoint
	return c
}

func (c *AuthenticatorBuilder) WithWellKnown(endpoint string) *AuthenticatorBuilder {
	c.wellKnownEndpoint = endpoint
	return c
}

func (c *AuthenticatorBuilder) WithRedirectUrl(url string) *AuthenticatorBuilder {
	c.redirectUrl = url
	return c
}

func (c *AuthenticatorBuilder) WithCertFile(filename string) *AuthenticatorBuilder {
	c.certFile = filename
	return c
}

func (c *AuthenticatorBuilder) WithKeyFile(filename string) *AuthenticatorBuilder {
	c.keyFile = filename
	return c
}

func (c *AuthenticatorBuilder) WithRootCAs(rootCAs []string) *AuthenticatorBuilder {
	c.rootCAs = rootCAs
	return c
}

func (c *AuthenticatorBuilder) makeSecuredTransport() Transport {
	return NewSecureTransport(
		c.certFile,
		c.keyFile,
		c.rootCAs,
	)
}

func (c *AuthenticatorBuilder) makeCredentialsGranter(config Configuration) CredentialsGranter {
	return NewCredentialGrander(
		c.makeSecuredTransport(),
		config.TokenEndpoint,
		c.client,
	)
}

func (c *AuthenticatorBuilder) makeAccessConsenter() AccessConsenter {
	return NewAccessConsenter(
		c.makeSecuredTransport(),
		c.accessConsentEndpoint,
		c.fapiFinancialId,
	)
}

func (c *AuthenticatorBuilder) makePSUAccessConsenter() PSUAccessConsenter {
	return NewPSUAccessConsenter(
		c.makeSecuredTransport(),
		c.accessConsentEndpoint,
		c.redirectUrl,
		c.client,
	)
}

func (c *AuthenticatorBuilder) makeTokenGenerator(config Configuration) TokenGenerator {
	return NewTokenGenerator(
		c.makeSecuredTransport(),
		config.TokenEndpoint,
		c.redirectUrl,
		c.client,
	)
}
