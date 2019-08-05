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
	cmd.Flags().StringSliceVar(&cc.credhubConfig.credPermissions, "credhub-permissions", strings.Split(os.Getenv("CREDHUB_PERMISSIONS"), ","), "CredHub permissions to add to new UAA cc")

	return cmd
}

func (cc *clientCreateCmd) run() error {

	// construct the API, and validate it
	api := uaa.New(cc.uaaConfig.endpoint, "").WithClientCredentials(cc.uaaConfig.adminClientIdentity, cc.uaaConfig.adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)
	err := api.Validate()
	if err != nil {
		cc.log.Info(err.Error())
	}

	client := uaa.Client{
		ClientID:             cc.targetClientIdentity,
		ClientSecret:         cc.newClientConfig.clientPwd,
		AccessTokenValidity:  cc.newClientConfig.clientTokenValidity,
		AuthorizedGrantTypes: cc.newClientConfig.clientGrantTypes,
		Scope:                cc.newClientConfig.clientScopes,
		Authorities:          cc.newClientConfig.clientAuthorities,
	}

	newClient, err := api.CreateClient(client)

	if err != nil {
		cc.log.Info("Failed to create UAA Client: " + err.Error())
		return err
	}

	cc.log.Info(newClient.DisplayName)

	if cc.credhubConfig.endpoint != "" && cc.credhubConfig.clientID != "" && cc.credhubConfig.clientPwd != "" && cc.credhubConfig.credPermissions != nil && cc.credhubConfig.credPath != "" {
		chAdmin, err := credhub.New(cc.credhubConfig.endpoint,
			credhub.SkipTLSValidation(true),
			credhub.Auth(auth.UaaClientCredentials(cc.credhubConfig.clientID, cc.credhubConfig.clientPwd)),
			credhub.AuthURL(cc.uaaConfig.endpoint),
		)

		if err != nil {
			cc.log.Error("Failed to connect to CredHub", err)
			return err
		}

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
