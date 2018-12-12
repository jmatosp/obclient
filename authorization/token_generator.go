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

type TokenGenerator interface {
	Request(Code) (Token, error)
}

type tokenGenerator struct {
	transport   Transport
	endpoint    string
	redirectUrl string
	client      Client
}

func NewTokenGenerator(transport Transport, endpoint string, redirectUrl string, client Client) TokenGenerator {
	return tokenGenerator{
		transport:   transport,
		endpoint:    endpoint,
		redirectUrl: redirectUrl,
		client:      client,
	}
}

func (t tokenGenerator) Request(code Code) (Token, error) {
	client, err := t.transport.Client()
	if err != nil {
		return NoToken, errors.Wrap(err, "error getting access token")
	}

	request, err := http.NewRequest(http.MethodPost, t.endpoint, t.authCodeGrantReader(code))
	if err != nil {
		return NoToken, errors.Wrap(err, "error getting access token")
	}
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", t.client.AuthHeader())

	response, err := client.Do(request)
	if err != nil {
		return NoToken, errors.Wrap(err, "error getting access token")
	}

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))

		return NoToken, errors.New("error getting access token unexpected status code")
	}

	var accessTokenResponse AccessTokenResponse
	if err = json.NewDecoder(response.Body).Decode(&accessTokenResponse); err != nil {
		return NoToken, errors.Wrap(err, "error getting access token")
	}

	return Token(accessTokenResponse), nil
}

type Token AccessTokenResponse

var NoToken = Token{}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	Id          string `json:"id_token"`
}

func (t tokenGenerator) authCodeGrantReader(code Code) io.Reader {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("scope", "accounts")
	data.Set("code", code.Value)
	data.Set("redirect_uri", t.redirectUrl)
	return strings.NewReader(data.Encode())
}
