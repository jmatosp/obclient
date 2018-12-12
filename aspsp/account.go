package aspsp

import (
	"github.com/pkg/errors"
	"time"
)

type Account interface {
	Id() AccountId
	Currency() string
	Type() string
	Subtype() string
	Nickname() string
	AccountIdentity() AccountIdentity
	Transactions(from, to time.Time) ([]Transaction, error)
}

type AccountId string

type account struct {
	id                AccountId
	currency          string
	accType           string
	subtype           string
	nickname          string
	identity          AccountIdentity
	transactionLoader TransactionLoaderFunc
}

func NewAccount(id AccountId, currency, accType, subtype, nickname string, identity AccountIdentity, transactionLoader TransactionLoaderFunc) Account {
	return account{
		id:                id,
		currency:          currency,
		accType:           accType,
		subtype:           subtype,
		nickname:          nickname,
		identity:          identity,
		transactionLoader: transactionLoader,
	}
}

func (a account) Id() AccountId {
	return a.id
}

func (a account) Currency() string {
	return a.currency
}

func (a account) Type() string {
	return a.accType
}

func (a account) Subtype() string {
	return a.subtype
}

func (a account) Nickname() string {
	return a.nickname
}

func (a account) AccountIdentity() AccountIdentity {
	return nil
}

func (a account) Transactions(from, to time.Time) ([]Transaction, error) {
	return a.transactionLoader(from, to)
}

type AccountIdentity interface {
	SchemaName() string
	Identification() string
	Name() string
	SecondaryIdentification() string
	Servicer() string
}

type accountIdentity struct {
	schemaName              string
	identification          string
	name                    string
	secondaryIdentification string
	servicer                string
}

func NewAccountIdentity(schemaName, identification, name, secondaryIdentification, servicer string) AccountIdentity {
	return accountIdentity{
		schemaName:              schemaName,
		identification:          identification,
		name:                    name,
		secondaryIdentification: secondaryIdentification,
		servicer:                servicer,
	}
}

func (ai accountIdentity) SchemaName() string {
	return ai.schemaName
}

func (ai accountIdentity) Identification() string {
	return ai.identification
}

func (ai accountIdentity) Name() string {
	return ai.name
}

func (ai accountIdentity) SecondaryIdentification() string {
	return ai.secondaryIdentification
}

func (ai accountIdentity) Servicer() string {
	return ai.servicer
}

type TransactionLoaderFunc func(from, to time.Time) ([]Transaction, error)

func ErrNotImplementedTransactionLoaderFunc(_, _ time.Time) ([]Transaction, error) {
	return nil, errors.New("transactions loader not implemented")
}

type Transaction struct {
	Id            string
	Reference     string
	Value         string
	Currency      string
	Credit        bool
	Status        string
	ValueDateTime time.Time
	Information   string
}
