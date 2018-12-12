package aspsp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmatosp/obclient/authorization"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type AccountLister interface {
	List() ([]Account, error)
}

type accountLister struct {
	transport       authorization.Transport
	endpoint        string
	fapiFinancialId string
	accessToken     authorization.Token
}

func NewAccountLister(transport authorization.Transport, endpoint, fapiFinancialId string, accessToken authorization.Token) AccountLister {
	return &accountLister{
		transport:       transport,
		endpoint:        endpoint,
		fapiFinancialId: fapiFinancialId,
		accessToken:     accessToken,
	}
}

func (a *accountLister) List() ([]Account, error) {
	client, err := a.transport.Client()
	if err != nil {
		return []Account{}, errors.Wrap(err, "error listing accounts")
	}

	request, err := http.NewRequest(http.MethodGet, a.endpoint+"/accounts", nil)
	if err != nil {
		return []Account{}, errors.Wrap(err, "error listing accounts")
	}
	request.Header.Set("Authorization", "Bearer "+a.accessToken.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-fapi-financial-id", a.fapiFinancialId)
	request.Header.Set("x-fapi-interaction-id", uuid.New().String())

	response, err := client.Do(request)
	if err != nil {
		return []Account{}, errors.Wrap(err, "error listing accounts")
	}

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		return []Account{}, errors.Wrap(err, "error listing accounts")
	}

	var accountsResponse AccountsResponse
	if err = json.NewDecoder(response.Body).Decode(&accountsResponse); err != nil {
		return []Account{}, errors.Wrap(err, "error listing accounts")
	}

	var accounts []Account
	for _, account := range accountsResponse.Data.Account {
		accounts = append(accounts, mapAccounts(account))
	}

	return accounts, nil
}

type AccountsResponse struct {
	Data AccountsDataResponse `json:"data"`
}

type AccountsDataResponse struct {
	Account []AccountsDataAccountResponse `json:"account"`
}

type AccountsDataAccountResponse struct {
	AccountId      string `json:"AccountId"`
	Currency       string `json:"Currency"`
	Nickname       string `json:"Nickname"`
	AccountType    string `json:"AccountType"`
	AccountSubType string `json:"AccountSubType"`
}

func mapAccounts(account AccountsDataAccountResponse) Account {
	return NewAccount(
		AccountId(account.AccountId),
		account.Currency,
		account.AccountType,
		account.AccountSubType,
		account.Nickname,
		NewAccountIdentity("", "", "", "", ""),
		ErrNotImplementedTransactionLoaderFunc,
	)
}
