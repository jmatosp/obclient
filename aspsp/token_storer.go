package aspsp

import (
	"encoding/json"
	"github.com/jmatosp/obclient/authorization"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type TokenStorer interface {
	Store(authorization.Token) error
	Get() (authorization.Token, error)
}

type fileTokenStorer struct {
	folder string
}

func NewFileTokenStorer(folder string) TokenStorer {
	return fileTokenStorer{
		folder: folder,
	}
}

func (s fileTokenStorer) Store(token authorization.Token) error {
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "error storing token")
	}

	err = ioutil.WriteFile(s.filename(), tokenJson, 0644)
	if err != nil {
		return errors.Wrap(err, "error storing token")
	}

	return nil
}

func (s fileTokenStorer) Get() (authorization.Token, error) {
	clientJson, err := ioutil.ReadFile(s.filename())
	if os.IsNotExist(err) {
		return authorization.NoToken, ErrNotFound
	} else if err != nil {
		return authorization.NoToken, errors.Wrap(err, "error getting token")
	}

	var token authorization.Token
	err = json.Unmarshal(clientJson, &token)
	if err != nil {
		return authorization.NoToken, errors.Wrap(err, "error getting token")
	}

	return token, nil
}

func (s fileTokenStorer) filename() string {
	return path.Join(s.folder, "token.json")
}
