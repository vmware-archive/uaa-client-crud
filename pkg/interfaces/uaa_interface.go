package interfaces

import (
	"github.com/cloudfoundry-community/go-uaa"
)

type UaaApi struct {
	UaaApi *uaa.API
}

//go:generate counterfeiter ./ UaaAPI
type UaaAPI interface {
	GetClient(clientID string) (*uaa.Client, error)
	CreateClient(client uaa.Client) (*uaa.Client, error)
	ChangeClientSecret(id string, newSecret string) error
	UpdateClient(client uaa.Client) (*uaa.Client, error)
	DeleteClient(clientID string) (*uaa.Client, error)
}

func NewUaaApi(target string, zoneID string, adminClientIdentity string, adminClientSecret string) (UaaAPI, error) {
	newAPI, err := uaa.New(target, uaa.WithClientCredentials(adminClientIdentity, adminClientSecret, uaa.JSONWebToken))
	if err != nil {
		return nil, err
	}
	return &UaaApi{newAPI}, nil
}

func (u *UaaApi) GetClient(clientId string) (*uaa.Client, error) {
	return u.UaaApi.GetClient(clientId)
}

func (u *UaaApi) CreateClient(client uaa.Client) (*uaa.Client, error) {
	return u.UaaApi.CreateClient(client)
}

func (u *UaaApi) ChangeClientSecret(id string, newSecret string) error {
	return u.UaaApi.ChangeClientSecret(id, newSecret)
}

func (u *UaaApi) UpdateClient(client uaa.Client) (*uaa.Client, error) {
	return u.UaaApi.UpdateClient(client)
}

func (u *UaaApi) DeleteClient(clientID string) (*uaa.Client, error) {
	return u.UaaApi.DeleteClient(clientID)
}
