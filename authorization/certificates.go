package authorization

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Certificate interface {
	PublicKey() (*rsa.PublicKey, error)
	PrivateKey() (*rsa.PrivateKey, error)
}

// safeCertificates provide a lazy load and minimal memory trace of cert itself
type safeCertificates struct {
	publicCertFilename  string
	privateCertFilename string
}

func NewSafeCertificates(publicCertFilename, privateCertFilename string) Certificate {
	return safeCertificates{
		publicCertFilename:  publicCertFilename,
		privateCertFilename: privateCertFilename,
	}
}

func (c safeCertificates) PublicKey() (*rsa.PublicKey, error) {
	fileContents, err := ioutil.ReadFile(c.publicCertFilename)
	if err != nil {
		return nil, errors.Wrap(err, "error loading public cert file")
	}

	publicCertificate, err := jwt.ParseRSAPublicKeyFromPEM(fileContents)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing public cert file")
	}

	return publicCertificate, nil
}

func (c safeCertificates) PrivateKey() (*rsa.PrivateKey, error) {
	fileContents, err := ioutil.ReadFile(c.privateCertFilename)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private cert file")
	}

	privateCertificate, err := jwt.ParseRSAPrivateKeyFromPEM(fileContents)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing private cert file")
	}

	err = privateCertificate.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "error validating private key")
	}

	return privateCertificate, err
}
