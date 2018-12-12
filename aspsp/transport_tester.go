package aspsp

import (
	"github.com/jmatosp/ob_security/authorization"
	"github.com/pkg/errors"
)

// TransportTester tests if the certificates are valid making a token call
type TransportTester interface {
	Test() error
}

type transportTester struct {
	authorization.Transport
	tokenEndpoint string
}

func NewTransportTester(transport authorization.Transport, tokenEndpoint string) TransportTester {
	return transportTester{transport, tokenEndpoint}
}

func (t transportTester) Test() error {
	client, err := t.Transport.Client()
	if err != nil {
		return errors.Wrap(err, "error testing transport")
	}

	_, err = client.Get(t.tokenEndpoint)
	if err != nil {
		return errors.Wrap(err, "error testing transport")
	}

	return nil
}
