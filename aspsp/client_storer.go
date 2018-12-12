package aspsp

import (
	"encoding/json"
	"github.com/jmatosp/ob_security/authorization"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

var ErrNotFound = errors.New("client not found")

type ClientStorer interface {
	Store(authorization.Client) error
	Get() (authorization.Client, error)
}

type fileStorer struct {
	folder string
}

func NewClientStorer(folder string) ClientStorer {
	return &fileStorer{
		folder: folder,
	}
}

func (s *fileStorer) Store(client authorization.Client) error {
	clientJson, err := json.Marshal(client)
	if err != nil {
		return errors.Wrap(err, "error storing client")
	}

	err = ioutil.WriteFile(s.filename(), clientJson, 0644)
	if err != nil {
		return errors.Wrap(err, "error storing client")
	}

	return nil
}

func (s *fileStorer) Get() (authorization.Client, error) {
	clientJson, err := ioutil.ReadFile(s.filename())
	if os.IsNotExist(err) {
		return authorization.NoClient, ErrNotFound
	} else if err != nil {
		return authorization.NoClient, errors.Wrap(err, "error getting client")
	}

	var client authorization.Client
	err = json.Unmarshal(clientJson, &client)
	if err != nil {
		return authorization.NoClient, errors.Wrap(err, "error getting client")
	}

	return client, nil
}

func (s *fileStorer) filename() string {
	return path.Join(s.folder, "client.json")
}
