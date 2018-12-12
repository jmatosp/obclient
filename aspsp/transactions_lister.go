package aspsp

import "time"

type TransactionLister interface {
	List(to, from *time.Time) ([]Transaction, error)
}
