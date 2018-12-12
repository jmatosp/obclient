package authorization

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type ClientRegister interface {
	Register() (Client, error)
}

func NewClientRegisterer(registrationEndpoint, issuer string, softwareStatement SoftwareStatement, transport Transport) ClientRegister {
	return &clientRegister{
		registrationEndpoint: registrationEndpoint,
		issuer:               issuer,
		softwareStatement:    softwareStatement,
		transport:            transport,
	}
}

type clientRegister struct {
	registrationEndpoint string
	issuer               string
	softwareStatement    SoftwareStatement
	transport            Transport
}

func (o *clientRegister) Register() (Client, error) {
	client, err := o.transport.Client()
	if err != nil {
		return NoClient, err
	}

	payload, err := o.signedRegisterClaims()
	if err != nil {
		return NoClient, errors.Wrap(err, "error registering client")
	}

	request, err := http.NewRequest(http.MethodPost, o.registrationEndpoint, bytes.NewBufferString(payload))
	if err != nil {
		return NoClient, errors.Wrap(err, "error registering client")
	}
	request.Header.Add("Content-Type", "application/jwt")

	response, err := client.Do(request)
	if err != nil {
		return NoClient, errors.Wrap(err, "error registering client")
	}

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		message, err := ioutil.ReadAll(response.Body)
		if err != nil {
			message = []byte("can't decode response body")
		}
		return NoClient, errors.Errorf("unexpected status code %d creating client: %s", response.StatusCode, string(message))
	}

	var registrationResponse OBClientRegistrationResponse
	if err = json.NewDecoder(response.Body).Decode(&registrationResponse); err != nil {
		return NoClient, errors.Wrap(err, "error registering client")
	}

	return mapToClient(registrationResponse), nil
}

func mapToClient(registrationResponse OBClientRegistrationResponse) Client {
	return NewClient(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
	)
}

func (o *clientRegister) signedRegisterClaims() (string, error) {
	claims := o.registerClaims()
	err := claims.Valid()
	if err != nil {
		return "", err
	}
	return o.softwareStatement.Sign(claims)
}

func (o *clientRegister) registerClaims() jwt.Claims {
	iat := time.Now()
	exp := iat.Add(time.Hour)
	return jwt.MapClaims{
		"kid":                             "YqL1S1MVsiknkoNpAMcXXui0VOQ",
		"token_endpoint_auth_signing_alg": "RS256",
		"grant_types": []string{
			"authorization_code",
			"refresh_token",
			"client_credentials",
		},
		"subject_type":     "public",
		"application_type": "web",
		"iss":              o.softwareStatement.Id(),
		"redirect_uris": []string{
			o.softwareStatement.RedirectUrl(),
		},
		"token_endpoint_auth_method": "client_secret_basic",
		"aud":                        o.issuer,
		"software_statement":         o.softwareStatement.Name(),
		"scopes": []string{
			"openid",
			"accounts",
			"payments",
		},
		"request_object_signing_alg": "none",
		"exp":                        exp.Unix(),
		"iat":                        iat.Unix(),
		"jti":                        uuid.New().String(),
		"response_types": []string{
			"code",
			"code id_token",
		},
		"id_token_signed_response_alg": "RS256",
	}
}

type OBClientRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
