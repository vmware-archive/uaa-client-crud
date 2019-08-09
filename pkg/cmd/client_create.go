package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

type clientCreateCmd struct {
	baseCmd
	newClientConfig uaaClientConfig
}

type uaaClientConfig struct {
	clientPwd        string
	clientGrantTypes []string

	clientScopes        []string
	clientAuthorities   []string
	clientTokenValidity int64
}

func NewCreateClientCmd(uaaApiFactory uaaApiFactory, credHubFactory credHubFactory, out io.Writer) *cobra.Command {
	cc := &clientCreateCmd{
		newBaseCmd(
			uaaApiFactory,
			credHubFactory, out),
		uaaClientConfig{},
	}

	cmd := &cobra.Command{
		Use:    "create",
		Short:  "Create a new client in UAA",
		PreRun: cc.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			trimWhitespace(cc.newClientConfig.clientGrantTypes)
			trimWhitespace(cc.newClientConfig.clientScopes)
			trimWhitespace(cc.newClientConfig.clientAuthorities)
			return cc.run()
		},
	}

	cc.addCommonFlags(cmd)

	cmd.Flags().StringVarP(&cc.newClientConfig.clientPwd, "target-client-pwd", "w", "", "Target Client Password")
	cmd.MarkFlagRequired("target-client-pwd")
	cmd.Flags().StringSliceVarP(&cc.newClientConfig.clientGrantTypes, "auth-grant-types", "g", []string{"client_credentials"}, "A comma separated list of Authorization Grant Types")
	cmd.Flags().StringSliceVarP(&cc.newClientConfig.clientScopes, "scopes", "s", []string{"uaa.none"}, "A comma separated list of UAA Scopes")
	cmd.Flags().StringSliceVarP(&cc.newClientConfig.clientAuthorities, "authorities", "a", []string{""}, "A comma separated list of UAA Authorities")
	cmd.MarkFlagRequired("authorities")
	cmd.Flags().Int64VarP(&cc.newClientConfig.clientTokenValidity, "token-validity", "t", 1800, "Client token validity period in seconds")
	cmd.Flags().StringSliceVar(&cc.credhubConfig.credPermissions, "credhub-permissions", strings.Split(strings.ReplaceAll(os.Getenv("CREDHUB_PERMISSIONS"), " ", ""), ","), "CredHub permissions to add to new UAA cc")

	return cmd
}

func (cc *clientCreateCmd) run() error {
	// construct the API, and validate it
	apiClient := cc.uaaApiFactory(cc.uaaConfig.endpoint, "", cc.uaaConfig.adminClientIdentity, cc.uaaConfig.adminClientPwd)

	err := apiClient.Validate()
	if err != nil {
		cc.log.Error("Error validating UUA API client", err)
		return err
	}

	cc.log.Info("connected to UAA")

	client := uaa.Client{
		ClientID:             cc.targetClientIdentity,
		ClientSecret:         cc.newClientConfig.clientPwd,
		AccessTokenValidity:  cc.newClientConfig.clientTokenValidity,
		AuthorizedGrantTypes: cc.newClientConfig.clientGrantTypes,
		Scope:                cc.newClientConfig.clientScopes,
		Authorities:          cc.newClientConfig.clientAuthorities,
	}
	c, err := apiClient.GetClient(client.ClientID)

	if err != nil {
		cc.log.Debug("UAA client does not exist. Creating.")
		newClient, err := apiClient.CreateClient(client)
		if err != nil {
			cc.log.Error("Failed to create UAA Client", err)
			return err
		} else {
			cc.log.Debug("Created UAA client [" + newClient.ClientID + "]")
		}
	} else {
		if c.ClientID == client.ClientID {
			cc.log.Info("Found existing client ID in UAA, updating")
			err := apiClient.ChangeClientSecret(client.ClientID, client.ClientSecret)
			if err != nil {
				cc.log.Error("Failed to update client secret", err)
				return err
			}
			c, err = apiClient.UpdateClient(client)
			if err != nil {
				cc.log.Error("Failed to update client", err)
				return err
			} else {
				cc.log.Debug("UAA client updated")
			}
		}
	}

	if cc.credhubConfig.endpoint != "" && cc.credhubConfig.clientID != "" && cc.credhubConfig.clientPwd != "" && cc.credhubConfig.credPermissions != nil && cc.credhubConfig.credPath != "" {
		cc.log.Debug("Found CredHub config")
		chAdmin, err := cc.credHubFactory(cc.credhubConfig.endpoint,
			true,
			cc.credhubConfig.clientID,
			cc.credhubConfig.clientPwd,
			cc.uaaConfig.endpoint,
		)

		if err != nil {
			cc.log.Error("Failed to connect to CredHub", err)
			return err
		}

		p, err := chAdmin.GetPermissionByPathActor(cc.credhubConfig.credPath, "uaa-client:"+cc.targetClientIdentity)
		if err != nil {
			_, err = chAdmin.AddPermission(
				cc.credhubConfig.credPath, "uaa-client:"+cc.targetClientIdentity,
				cc.credhubConfig.credPermissions,
			)
			if err != nil {
				cc.log.Error("Failed to add CredHub permission", err)
				return err
			} else {
				cc.log.Debug("CredHub permission created")
			}
		} else {
			_, err = chAdmin.UpdatePermission(
				p.UUID,
				p.Path, p.Actor,
				cc.credhubConfig.credPermissions,
			)
			if err != nil {
				cc.log.Error("Failed to update CredHub permission", err)
				return err
			} else {
				cc.log.Debug("CredHub permission updated")
			}

		}

	}
	return nil
}

func trimWhitespace(a []string) {
	for i, v := range a {
		a[i] = strings.TrimSpace(v)
	}
}
