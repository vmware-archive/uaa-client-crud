package cmd

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

type clientDeleteCmd struct {
	baseCmd
}

func NewDeleteClientCmd(out io.Writer) *cobra.Command {
	client := &clientDeleteCmd{
		baseCmd: newBaseCmd(out),
	}

	cd := &clientDeleteCmd{}
	cmd := &cobra.Command{
		Use:    "create",
		Short:  "Create a new client in UAA",
		PreRun: cd.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return cd.run()
		},
	}

	cd.addCommonFlags(cmd)
	cmd.Flags().StringVarP(&client.uaaConfig.endpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cobra.MarkFlagRequired(cmd.Flags(), "uaa-endpoint")
	cmd.Flags().StringVarP(&client.uaaConfig.adminClientIdentity, "admin-identity", "i", "", "Admin Username")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-identity")
	cmd.Flags().StringVarP(&client.uaaConfig.adminClientPwd, "admin-pwd", "p", "", "Admin Password")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-pwd")
	cmd.Flags().StringVarP(&client.newClientConfig.clientIndentity, "client-identity", "c", "", "New Client Identity")
	cobra.MarkFlagRequired(cmd.Flags(), "client-identity")
	cmd.Flags().StringVarP(&client.newClientConfig.clientPwd, "client-pwd", "w", "", "New Client Password")
	cobra.MarkFlagRequired(cmd.Flags(), "client-pwd")
	cmd.Flags().StringSliceVarP(&client.newClientConfig.clientGrantTypes, "auth-grant-types", "g", []string{"client_credentials"}, "A comma separated list of Authorization Grant Types")
	cmd.Flags().StringSliceVarP(&client.newClientConfig.clientScopes, "scopes", "s", []string{"uaa.none"}, "A comma separated list of UAA Scopes")
	cmd.Flags().StringSliceVarP(&client.newClientConfig.clientAuthorities, "authorities", "a", []string{""}, "A comma separated list of UAA Authorities")
	cobra.MarkFlagRequired(cmd.Flags(), "authorities")
	cmd.Flags().Int64VarP(&client.newClientConfig.clientTokenValidity, "token-validity", "t", 1800, "Client token validity period in seconds")
	cmd.Flags().StringVar(&client.credhubConfig.clientID, "credhub-client-identity", os.Getenv("CREDHUB_CLIENT_ID"), "Credhub Client Identity if granting new client CredHub access")
	cmd.Flags().StringVar(&client.credhubConfig.endpoint, "credhub-endpoint", os.Getenv("CREDHUB_URL"), "CredHub endpoint URL")
	cmd.Flags().StringVar(&client.credhubConfig.credPath, "credential-path", os.Getenv("CREDHUB_CRED_PATH"), "CredHub Credential Path")
	cmd.Flags().StringSliceVar(&client.credhubConfig.credPermissions, "credhub-permissions", strings.Split(os.Getenv("CREDHUB_PERMISSIONS"), ","), "CredHub permissions to add to new UAA client")

	return cmd
}

func (cd *clientDeleteCmd) run() error {

	// construct the API, and validate it
	api := uaa.New(cd.uaaConfig.endpoint, "").WithClientCredentials(cd.uaaConfig.adminClientIdentity, cd.uaaConfig.adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)
	err := api.Validate()
	if err != nil {
		cd.log.Error("", err)
	}

	_, err = api.DeleteClient(cd.newClientConfig.clientIndentity)

	if err != nil {
		cd.log.Error("Failed to delete UAA Client", err)
		return err
	}

	if cd.credhubConfig.endpoint != "" && cd.credhubConfig.clientID != "" && cd.credhubConfig.credPermissions != nil && cd.credhubConfig.credPath != "" {

		chAdmin, err := credhub.New(cd.credhubConfig.endpoint,
			credhub.SkipTLSValidation(true),
			credhub.Auth(auth.UaaClientCredentials(cd.credhubConfig.clientID, cd.credhubConfig.clientPwd)),
			credhub.AuthURL(cd.uaaConfig.endpoint),
		)

		if err != nil {
			cd.log.Error("Failed to connect to CredHub", err)
			return err
		}

		permission, err := chAdmin.GetPermissionByPathActor(cd.credhubConfig.credPath, "uaa-client:"+cd.newClientConfig.clientIndentity)
		if err != nil {
			cd.log.Error("Failed to get Permission object", err)
			return err
		}
		_, err = chAdmin.DeletePermission(permission.UUID)
		if err != nil {
			cd.log.Error("Failed to delete Permission object", err)
			return err
		}

	}
	return nil
}
