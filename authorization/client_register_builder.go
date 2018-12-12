package authorization

import (
	"github.com/pkg/errors"
)

type ClientRegisterBuilder struct {
	wellKnownEndpoint     string
	sigPublicKeyFile      string
	sigPrivateKeyFile     string
	softwareStatementID   string
	softwareStatementName string
	redirectUrl           string
	certFile              string
	keyFile               string
	rootCAs               []string
}

func NewClientRegisterBuilder() *ClientRegisterBuilder {
	return &ClientRegisterBuilder{}
}

func (c *ClientRegisterBuilder) Build() (ClientRegister, error) {
	err := c.mustValidate()

	config, err := GetConfiguration(c.wellKnownEndpoint)
	if err != nil {
		return nil, err
	}

	return NewClientRegisterer(
		config.RegistrationEndpoint,
		config.Issuer,
		c.makeSoftwareStatement(),
		c.makeSecuredTransport(),
	), nil
}

func (c *ClientRegisterBuilder) mustValidate() error {
	if c.wellKnownEndpoint == "" {
		return errors.New("error wellKnownEndpoint not provided")
	}

	if c.sigPublicKeyFile == "" {
		return errors.New("error sigPublicKeyFile not provided")
	}

	if c.sigPrivateKeyFile == "" {
		return errors.New("error sigPrivateKeyFile not provided")
	}

	if c.softwareStatementID == "" {
		return errors.New("error softwareStatementID not provided")
	}

	if c.softwareStatementName == "" {
		return errors.New("error softwareStatementName not provided")
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

func (c *ClientRegisterBuilder) WithWellKnown(endpoint string) *ClientRegisterBuilder {
	c.wellKnownEndpoint = endpoint
	return c
}

func (c *ClientRegisterBuilder) WithSigPublicKeyFile(filename string) *ClientRegisterBuilder {
	c.sigPublicKeyFile = filename
	return c
}

func (c *ClientRegisterBuilder) WithSigPrivateKeyFile(filename string) *ClientRegisterBuilder {
	c.sigPrivateKeyFile = filename
	return c
}

func (c *ClientRegisterBuilder) WithSoftwareStatementID(id string) *ClientRegisterBuilder {
	c.softwareStatementID = id
	return c
}

func (c *ClientRegisterBuilder) WithSoftwareStatementName(name string) *ClientRegisterBuilder {
	c.softwareStatementName = name
	return c
}

func (c *ClientRegisterBuilder) WithRedirectUrl(url string) *ClientRegisterBuilder {
	c.redirectUrl = url
	return c
}

func (c *ClientRegisterBuilder) WithCertFile(filename string) *ClientRegisterBuilder {
	c.certFile = filename
	return c
}

func (c *ClientRegisterBuilder) WithKeyFile(filename string) *ClientRegisterBuilder {
	c.keyFile = filename
	return c
}

func (c *ClientRegisterBuilder) WithRootCAs(rootCAs []string) *ClientRegisterBuilder {
	c.rootCAs = rootCAs
	return c
}

func (c *ClientRegisterBuilder) makeSoftwareStatement() SoftwareStatement {
	signer := NewSafeCertificates(
		c.sigPublicKeyFile,
		c.sigPrivateKeyFile,
	)

	return NewSoftwareStatement(
		c.softwareStatementID,
		c.softwareStatementName,
		c.redirectUrl,
		NewRS256Singer(signer),
	)
}

func (c *ClientRegisterBuilder) makeSecuredTransport() Transport {
	return NewSecureTransport(
		c.certFile,
		c.keyFile,
		c.rootCAs,
	)
}
