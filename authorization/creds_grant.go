package authorization

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CredentialsGranter interface {
	Request() (GrantToken, error)
}

type credentialsGranter struct {
	transport Transport
	endpoint  string
	client    Client
}

func NewCredentialGrander(transport Transport, tokenEndpoint string, client Client) CredentialsGranter {
	return credentialsGranter{
		transport: transport,
		endpoint:  tokenEndpoint,
		client:    client,
	}
}

func (c credentialsGranter) Request() (GrantToken, error) {
	client, err := c.transport.Client()
	if err != nil {
		return NoGrantToken, errors.Wrap(err, "error getting credentials grant")
	}

	request, err := http.NewRequest(http.MethodPost, c.endpoint, credentialsGrantRequestReader())
	if err != nil {
		return NoGrantToken, errors.Wrap(err, "error getting credentials grant")
	}
	request.Header.Set("Authorization", c.client.AuthHeader())
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return NoGrantToken, errors.Wrap(err, "error getting credentials grant")
	}

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))

		return NoGrantToken, errors.New("error getting credentials grant unexpected status code")
	}

	var credentialsGrantResponse CredentialsGrantResponse
	if err = json.NewDecoder(response.Body).Decode(&credentialsGrantResponse); err != nil {
		return NoGrantToken, errors.Wrap(err, "error getting credentials grant")
	}

	return GrantToken(credentialsGrantResponse), nil
}

type GrantToken CredentialsGrantResponse

var NoGrantToken = GrantToken{}

type CredentialsGrantResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func credentialsGrantRequestReader() io.Reader {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "accounts openid")
	return strings.NewReader(data.Encode())
}
