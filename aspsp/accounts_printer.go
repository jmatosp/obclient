package aspsp

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type AccountsPrinter struct {
	w *tabwriter.Writer
}

func NewAccountsPrinter() AccountsPrinter {
	return AccountsPrinter{
		tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug),
	}
}

func (a AccountsPrinter) Print(accounts []Account) {
	a.header()
	for _, account := range accounts {
		a.accountPrint(account)
	}
	a.w.Flush()
}

func (a AccountsPrinter) header() {
	fmt.Fprintf(a.w, "Id\tCurrency\tNickname\tType\tSubType\n")
}

func (a AccountsPrinter) accountPrint(account Account) {
	fmt.Fprintf(a.w, "%s\t%s\t%s\t%s\t%s\n",
		account.Id(),
		account.Currency(),
		account.Nickname(),
		account.Type(),
		account.Subtype(),
	)
}
