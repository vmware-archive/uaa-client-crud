package interfaces

import (
	"github.com/cloudfoundry-community/go-uaa"
)

type UaaApi struct {
	UaaApi *uaa.API
}

type uaaAPI interface {
	Validate() error
	GetClient(clientID string) (*uaa.Client, error)
	CreateClient(client uaa.Client) (*uaa.Client, error)
	ChangeClientSecret(id string, newSecret string) error
}

func NewUaaApi(target string, zoneID string, adminClientIdentity string, adminClientPwd string) *UaaApi {
	return &UaaApi{uaa.New(target, zoneID).WithClientCredentials(adminClientIdentity, adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)}
}
