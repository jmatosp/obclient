package authorization

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Transport interface {
	Client() (*http.Client, error)
}

type secureTransport struct {
	cerFile string
	keyFile string
	certs   []string
	conn    *http.Client
}

func NewSecureTransport(cerFile, keyFile string, certs []string) Transport {
	return &secureTransport{
		cerFile: cerFile,
		keyFile: keyFile,
		certs:   certs,
	}
}

func (t *secureTransport) Client() (*http.Client, error) {
	var err error
	if t.conn == nil {
		t.conn, err = t.client()
	}
	return t.conn, err
}

func (t *secureTransport) client() (*http.Client, error) {
	pool := x509.NewCertPool()
	for _, cert := range t.certs {
		ca, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading cert %s", cert)
		}

		if ok := pool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.Wrapf(err, "error appending cert %s", cert)
		}
	}

	clientCert, err := tls.LoadX509KeyPair(t.cerFile, t.keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "error loading certFile and keyFile")
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      pool,
	}

	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}

	return &http.Client{
		Timeout:   time.Minute * 2,
		Transport: &transport,
	}, nil
}
