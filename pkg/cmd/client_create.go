package cmd

import (
	"io"
	"os"
	"strings"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

type clientCreateCmd struct {
	baseCmd
	newClientConfig uaaClientConfig
}

type uaaClientConfig struct {
	clientPwd           string
	clientGrantTypes    []string
	clientScopes        []string
	clientAuthorities   []string
	clientTokenValidity int64
}

func NewCreateClientCmd(out io.Writer) *cobra.Command {
	cc := &clientCreateCmd{
		newBaseCmd(out),
		uaaClientConfig{},
	}

	cmd := &cobra.Command{
		Use:    "create",
		Short:  "Create a new client in UAA",
		PreRun: cc.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
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
	trimWhitespace(cc.newClientConfig.clientGrantTypes)
	trimWhitespace(cc.newClientConfig.clientScopes)
	trimWhitespace(cc.newClientConfig.clientAuthorities)

	// construct the API, and validate it
	api := uaa.New(cc.uaaConfig.endpoint, "").WithClientCredentials(cc.uaaConfig.adminClientIdentity, cc.uaaConfig.adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)
	err := api.Validate()
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
	cc.log.Info("checking if client already exists")
	c, err := api.GetClient(client.ClientID)

	if err != nil {
		cc.log.Info("UAA client does not exist. Creating")
		newClient, err2 := api.CreateClient(client)
		if err2 != nil {
			cc.log.Error("Failed to create UAA Client", err)
			return err2
		}
		cc.log.Info(newClient.DisplayName)
	} else {
		if c.ClientID == client.ClientID {
			cc.log.Info("Found existing client ID in UAA, updating")
			err := api.ChangeClientSecret(client.ClientID, client.ClientSecret)
			if err != nil {
				cc.log.Error("Failed to update client secret", err)
				return err
			}
			c, err = api.UpdateClient(client)
			if err != nil {
				cc.log.Error("Failed to update client", err)
				return err
			}
		}
	}

	if cc.credhubConfig.endpoint != "" && cc.credhubConfig.clientID != "" && cc.credhubConfig.clientPwd != "" && cc.credhubConfig.credPermissions != nil && cc.credhubConfig.credPath != "" {
		cc.log.Info("connecting to credhub")
		chAdmin, err := credhub.New(cc.credhubConfig.endpoint,
			credhub.SkipTLSValidation(true),
			credhub.Auth(auth.UaaClientCredentials(cc.credhubConfig.clientID, cc.credhubConfig.clientPwd)),
			credhub.AuthURL(cc.uaaConfig.endpoint),
		)

		if err != nil {
			cc.log.Error("Failed to connect to CredHub", err)
			return err
		}
		cc.log.Info("adding permission in credhub")
		_, err = chAdmin.AddPermission(
			cc.credhubConfig.credPath, "uaa-client:"+cc.targetClientIdentity,
			cc.credhubConfig.credPermissions,
		)

		if err != nil {
			cc.log.Error("Failed to add CredHub permission", err)
			return err
		}
	}
	return nil
}

func trimWhitespace(a []string) {
	for i, v := range a {
		a[i] = strings.TrimSpace(v)
	}
}
