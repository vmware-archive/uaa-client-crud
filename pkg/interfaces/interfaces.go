package interfaces

import (
	"github.com/cloudfoundry-community/go-uaa"
)

type myUaaClient interface {
	New(target string, zoneID string)
}

type uaaAPI interface {
	Validate() error
	GetClient(clientID string) (*uaa.Client, error)
	CreateClient(client uaa.Client) (*uaa.Client, error)
	ChangeClientSecret(id string, newSecret string) error
}

type UaaApi struct {
	uaaApi uaaAPI
}
