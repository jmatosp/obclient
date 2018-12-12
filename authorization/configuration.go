package authorization

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type Configuration struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	RegistrationEndpoint  string `json:"registration_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	Issuer                string `json:"issuer"`
}

var NoConfiguration = Configuration{}

func GetConfiguration(endpoint string) (Configuration, error) {
	response, err := http.Get(endpoint)
	if err != nil {
		return NoConfiguration, errors.Wrap(err, "error getting openid configuration")
	}

	var configuration Configuration
	if err = json.NewDecoder(response.Body).Decode(&configuration); err != nil {
		return NoConfiguration, errors.Wrap(err, "error getting openid configuration")
	}

	return configuration, err
}
