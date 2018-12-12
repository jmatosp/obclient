package authorization

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"time"
)

type PSUAccessConsenter interface {
	Request(AccessConsent) (Code, error)
}

type psuAccessConsenter struct {
	transport    Transport
	endpoint     string
	authCallback string
	client       Client
}

func NewPSUAccessConsenter(transport Transport, endpoint, authCallback string, client Client) PSUAccessConsenter {
	return psuAccessConsenter{
		transport:    transport,
		endpoint:     endpoint,
		authCallback: authCallback,
		client:       client,
	}
}

func (a psuAccessConsenter) Request(accessConsent AccessConsent) (Code, error) {
	client, err := a.transport.Client()
	if err != nil {
		return NoCode, errors.Wrap(err, "error starting user access consent flow")
	}

	endpoint := a.endpoint
	endpoint = fmt.Sprintf("https://modelobank2018.o3bank.co.uk:4501/ozone/v1.0/auth-code-url/%s?scope=%s", accessConsent.ConsentId, "openid%20accounts")

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return NoCode, errors.Wrap(err, "error starting user access consent flow")
	}
	request.Header.Set("Authorization", a.client.AuthHeader())

	response, err := client.Do(request)
	if err != nil {
		return NoCode, errors.Wrap(err, "error starting user access consent flow")
	}

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		return NoCode, errors.Errorf("error starting user access consent flow: unexpected response status code %d", response.StatusCode)
	}

	codeChan := make(chan Code)
	a.runCallbackListener(codeChan)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return NoCode, errors.Wrap(err, "error starting user access consent flow")
	}

	err = open.Run(string(body))
	if err != nil {
		return NoCode, errors.Wrap(err, "error initiating browser for user consent flow")
	}

	code := <-codeChan

	return code, nil
}

func (a psuAccessConsenter) runCallbackListener(tokenChan chan Code) {
	srv := &http.Server{Addr: ":8081"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		tokenChan <- Code{code}
		w.Write([]byte(`
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 
Transitional//EN"> <HTML> <HEAD> 
<TITLE>Open Banking Access Consent</TITLE> </HEAD>
<BODY>
<table border="0" width="100%">
 <tr>
  <td align="center"><font color=#330066 size="4"><strong>
   Authenticated!</strong></font>
  </td>
 </tr>
 <tr>
  <td align="center"><font color=#330066>
   (Please close this window)</font>
  </td>
 </tr>
</table>
</BODY>
</HTML> 
`))
		go func() {
			time.Sleep(time.Second * 2)
			srv.Shutdown(context.Background())
		}()
	})

	go func() {
		srv.ListenAndServe()
	}()
}

var NoCode = Code{}

type Code struct {
	Value string
}
