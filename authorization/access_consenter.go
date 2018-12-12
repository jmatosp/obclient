package authorization

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type AccessConsenter interface {
	Request(GrantToken) (AccessConsent, error)
}

type accessConsenter struct {
	transport       Transport
	endpoint        string
	fapiFinancialId string
}

func NewAccessConsenter(transport Transport, endpoint, fapiFinancialId string) AccessConsenter {
	return accessConsenter{
		transport:       transport,
		endpoint:        endpoint,
		fapiFinancialId: fapiFinancialId,
	}
}

func (a accessConsenter) Request(token GrantToken) (AccessConsent, error) {
	client, err := a.transport.Client()
	if err != nil {
		return NoAccessConsent, errors.Wrap(err, "error getting access consent")
	}

	data, err := json.Marshal(AccountsReadConsent)
	if err != nil {
		return NoAccessConsent, errors.Wrap(err, "error getting access consent")
	}

	request, err := http.NewRequest(http.MethodPost, a.endpoint+"/account-access-consents", bytes.NewBuffer(data))
	if err != nil {
		return NoAccessConsent, errors.Wrap(err, "error getting access consent")
	}
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-fapi-financial-id", a.fapiFinancialId)
	request.Header.Set("x-fapi-interaction-id", uuid.New().String())

	response, err := client.Do(request)
	if err != nil {
		return NoAccessConsent, errors.Wrap(err, "error getting access consent")
	}

	if response.StatusCode != http.StatusCreated {
		return NoAccessConsent, errors.Errorf("error getting access consent: unexpected response status code %d", response.StatusCode)
	}

	var accessConsentResponse AccessConsentResponse
	if err = json.NewDecoder(response.Body).Decode(&accessConsentResponse); err != nil {
		return NoAccessConsent, errors.Wrap(err, "error getting access consent")
	}

	return AccessConsent{accessConsentResponse.Data.ConsentId}, nil
}

var NoAccessConsent = AccessConsent{}

type AccessConsent struct {
	ConsentId string
}

type AccessConsentResponse struct {
	Data AccessConsentDataResponse `json:"Data"`
}

type AccessConsentDataResponse struct {
	ConsentId string `json:"ConsentId"`
}

type AccessConsentRequest struct {
	Data AccessConsentDataRequest `json:"Data"`
	Risk map[string]string        `json:"Risk"`
}

type AccessConsentDataRequest struct {
	AccountRequestID        string   `json:"AccountRequestId,omitempty"`
	Status                  string   `json:"Status,omitempty"`
	CreationDateTime        string   `json:"CreationDateTime,omitempty"`
	Permissions             []string `json:"Permissions"`
	ExpirationDateTime      string   `json:"ExpirationDateTime,omitempty"`
	TransactionFromDateTime string   `json:"TransactionFromDateTime"`
	TransactionToDateTime   string   `json:"TransactionToDateTime"`
}

var oneYearDuration = time.Hour * 24 * 365

var AccountsReadConsent = AccessConsentRequest{
	Data: AccessConsentDataRequest{
		Permissions: []string{
			"ReadAccountsBasic",
			"ReadAccountsDetail",
			"ReadBalances",
			"ReadBeneficiariesDetail",
			"ReadDirectDebits",
			"ReadProducts",
			"ReadStandingOrdersDetail",
			"ReadTransactionsCredits",
			"ReadTransactionsDebits",
			"ReadTransactionsDetail",
		},
		TransactionFromDateTime: time.Now().Add(-1 * oneYearDuration).Format("2006-01-02T15:04:05+00:00"),
		TransactionToDateTime:   time.Now().Add(oneYearDuration).Format("2006-01-02T15:04:05+00:00"),
	},
	Risk: map[string]string{},
}
