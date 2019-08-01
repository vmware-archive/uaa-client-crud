package cmd

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewDeleteClientCmd(logger lager.Logger) *cobra.Command {
	client := &uaaClient{}
	client.logger = logger
	uaaCreateClientCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new client in UAA",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return client.uaaDelete()
		},
	}

	uaaCreateClientCmd.Flags().StringVarP(&client.uaaConfig.endpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "uaa-endpoint")

	uaaCreateClientCmd.Flags().StringVarP(&client.uaaConfig.adminClientIdentity, "admin-identity", "i", "", "Admin Username")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "admin-identity")

	uaaCreateClientCmd.Flags().StringVarP(&client.uaaConfig.adminClientPwd, "admin-pwd", "p", "", "Admin Password")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "admin-pwd")

	uaaCreateClientCmd.Flags().StringVarP(&client.newClientConfig.clientIndentity, "client-identity", "c", "", "New Client Identity")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "client-identity")

	uaaCreateClientCmd.Flags().StringVarP(&client.newClientConfig.clientPwd, "client-pwd", "w", "", "New Client Password")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "client-pwd")

	uaaCreateClientCmd.Flags().StringSliceVarP(&client.newClientConfig.clientGrantTypes, "auth-grant-types", "g", []string{"client_credentials"}, "A comma separated list of Authorization Grant Types")

	uaaCreateClientCmd.Flags().StringSliceVarP(&client.newClientConfig.clientScopes, "scopes", "s", []string{"uaa.none"}, "A comma separated list of UAA Scopes")

	uaaCreateClientCmd.Flags().StringSliceVarP(&client.newClientConfig.clientAuthorities, "authorities", "a", []string{""}, "A comma separated list of UAA Authorities")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "authorities")

	uaaCreateClientCmd.Flags().Int64VarP(&client.newClientConfig.clientTokenValidity, "token-validity", "t", 1800, "Client token validity period in seconds")

	uaaCreateClientCmd.Flags().StringVar(&client.credhubConfig.clientID, "credhub-client-identity", os.Getenv("CREDHUB_CLIENT_ID"), "Credhub Client Identity if granting new client Credhub access")

	uaaCreateClientCmd.Flags().StringVar(&client.credhubConfig.endpoint, "credhub-endpoint", os.Getenv("CREDHUB_URL"), "Credhub endpoint URL")

	uaaCreateClientCmd.Flags().StringVar(&client.credhubConfig.credPath, "credential-path", os.Getenv("CREDHUB_CRED_PATH"), "Credhub Credential Path")

	uaaCreateClientCmd.Flags().StringSliceVar(&client.credhubConfig.credPermissions, "credhub-permissions", strings.Split(os.Getenv("CREDHUB_PERMISSIONS"), ","), "Credhub permissions to add to new UAA client")

	return uaaCreateClientCmd
}

func (cc *uaaClient) uaaDelete() error {

	// construct the API, and validate it
	api := uaa.New(cc.uaaConfig.endpoint, "").WithClientCredentials(cc.uaaConfig.adminClientIdentity, cc.uaaConfig.adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)
	err := api.Validate()
	if err != nil {
		cc.logger.Info(err.Error())
	}

	_, err = api.DeleteClient(cc.newClientConfig.clientIndentity)

	if err != nil {
		cc.logger.Info("Failed to delete UAA Client: " + err.Error())
		return err
	}

	if cc.credhubConfig.endpoint != "" && cc.credhubConfig.clientID != "" && cc.credhubConfig.credPermissions != nil && cc.credhubConfig.credPath != "" {
		credHubClient, err := api.GetClient(cc.credhubConfig.clientID)
		cc.logger.Info(credHubClient.ClientID + " secret:" + credHubClient.ClientSecret)
		err = credHubClient.Validate()
		if err != nil {
			cc.logger.Info(err.Error())
			return err
		}

		chAdmin, err := credhub.New(cc.credhubConfig.endpoint,
			credhub.SkipTLSValidation(true),
			credhub.Auth(auth.UaaClientCredentials(cc.credhubConfig.clientID, "iFB7oFXyRI1Yp3sHd_5RZ7WLDZHv2UX3")),
			credhub.AuthURL(cc.uaaConfig.endpoint),
		)

		permission, err := chAdmin.GetPermissionByPathActor(cc.credhubConfig.credPath, "uaa-client:"+cc.newClientConfig.clientIndentity)

		_, err = chAdmin.DeletePermission(permission.UUID)

	}
	return nil
}
