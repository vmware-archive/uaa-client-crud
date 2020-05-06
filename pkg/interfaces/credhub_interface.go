package interfaces

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/permissions"
	"fmt"
)

type CredHubApi struct {
	CredHub *credhub.CredHub
}

//go:generate counterfeiter ./ CredHubAPI
type CredHubAPI interface {
	//methods to implement
	GetPermissionByPathActor(path string, actor string) (*permissions.Permission, error)
	AddPermission(path string, actor string, ops []string) (*permissions.Permission, error)
	UpdatePermission(uuid string, path string, actor string, ops []string) (*permissions.Permission, error)
	DeletePermission(uuid string) (*permissions.Permission, error)
	DeleteCredential(name string) error
	FindByPath(path string) ([]string, error)
	PrintCredhub() string
}

func NewCredHubApi(target string, skipTLS bool, clientID string, clientSecret string, uaaEndpoint string) (CredHubAPI, error) {

	api, err := credhub.New(target,
		credhub.SkipTLSValidation(skipTLS),
		credhub.Auth(auth.UaaClientCredentials(clientID, clientSecret)),
		credhub.AuthURL(uaaEndpoint))
	return &CredHubApi{api}, err
}

func (c *CredHubApi) GetPermissionByPathActor(path string, actor string) (*permissions.Permission, error) {
	return c.CredHub.GetPermissionByPathActor(path, actor)
}

func (c *CredHubApi) AddPermission(path string, actor string, ops []string) (*permissions.Permission, error) {
	return c.CredHub.AddPermission(path, actor, ops)
}

func (c *CredHubApi) UpdatePermission(uuid string, path string, actor string, ops []string) (*permissions.Permission, error) {
	return c.CredHub.UpdatePermission(uuid, path, actor, ops)
}

func (c *CredHubApi) DeletePermission(uuid string) (*permissions.Permission, error) {
	return c.CredHub.DeletePermission(uuid)
}

func (c *CredHubApi) DeleteCredential(name string) error {
	return c.CredHub.Delete(name)
}

func (c *CredHubApi) FindByPath(path string) ([]string, error) {
	results, err := c.CredHub.FindByPath(path)
	if err != nil {
		return nil, err
	}
	var credpaths []string
	for _, result := range results.Credentials {
		credpaths = append(credpaths, result.Name)
	}
	return credpaths, nil
}

func (c *CredHubApi) PrintCredhub() string {
	return fmt.Sprintf("Credhub info: [%+v]\n "+
		"Credhub auth: [%+v] \n"+
		"Credhub API URL: [%+v]", c.CredHub.Info(), c.CredHub.Auth, c.CredHub.ApiURL)
}
