package authorization

type SoftwareStatement interface {
	Signer
	Id() string
	Name() string
	RedirectUrl() string
}

type softwareStatement struct {
	Signer
	id          string
	name        string
	redirectUrl string
}

func NewSoftwareStatement(id, name, redirectUrl string, signer Signer) SoftwareStatement {
	return softwareStatement{
		Signer:      signer,
		id:          id,
		name:        name,
		redirectUrl: redirectUrl,
	}
}

func (s softwareStatement) Id() string {
	return s.id
}

func (s softwareStatement) Name() string {
	return s.name
}

func (s softwareStatement) RedirectUrl() string {
	return s.redirectUrl
}
