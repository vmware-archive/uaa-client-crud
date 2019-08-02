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

type clientCreateCmd struct {
	baseCmd
}

func NewCreateClientCmd(out io.Writer) *cobra.Command {
	cc := &clientCreateCmd{
		newBaseCmd(out),
	}

	cmd := &cobra.Command{
		Use:    "create",
		Short:  "Create a new cc in UAA",
		PreRun: cc.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return cc.run()
		},
	}

	cc.addCommonFlags(cmd)

	cmd.Flags().StringVarP(&cc.uaaConfig.endpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cobra.MarkFlagRequired(cmd.Flags(), "uaa-endpoint")
	cmd.Flags().StringVarP(&cc.uaaConfig.adminClientIdentity, "admin-identity", "i", "", "Admin Username")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-identity")
	cmd.Flags().StringVarP(&cc.uaaConfig.adminClientPwd, "admin-pwd", "p", "", "Admin Password")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-pwd")
	cmd.Flags().StringVarP(&cc.newClientConfig.clientIndentity, "cc-identity", "c", "", "New Client Identity")
	cobra.MarkFlagRequired(cmd.Flags(), "cc-identity")
	cmd.Flags().StringVarP(&cc.newClientConfig.clientPwd, "cc-pwd", "w", "", "New Client Password")
	cobra.MarkFlagRequired(cmd.Flags(), "cc-pwd")
	cmd.Flags().StringSliceVarP(&cc.newClientConfig.clientGrantTypes, "auth-grant-types", "g", []string{"client_credentials"}, "A comma separated list of Authorization Grant Types")
	cmd.Flags().StringSliceVarP(&cc.newClientConfig.clientScopes, "scopes", "s", []string{"uaa.none"}, "A comma separated list of UAA Scopes")
	cmd.Flags().StringSliceVarP(&cc.newClientConfig.clientAuthorities, "authorities", "a", []string{""}, "A comma separated list of UAA Authorities")
	cobra.MarkFlagRequired(cmd.Flags(), "authorities")
	cmd.Flags().Int64VarP(&cc.newClientConfig.clientTokenValidity, "token-validity", "t", 1800, "Client token validity period in seconds")
	cmd.Flags().StringVar(&cc.credhubConfig.clientID, "credhub-cc-identity", os.Getenv("CREDHUB_CLIENT_ID"), "CredHub Client Identity if granting new cc Credhub access")
	cmd.Flags().StringVar(&cc.credhubConfig.clientPwd, "credhub-cc-password", os.Getenv("CREDHUB_CLIENT_PASSWORD"), "Credhub Client Password if granting the new cc CredHub access")
	cmd.Flags().StringVar(&cc.credhubConfig.endpoint, "credhub-endpoint", os.Getenv("CREDHUB_URL"), "CredHub endpoint URL")
	cmd.Flags().StringVar(&cc.credhubConfig.credPath, "credential-path", os.Getenv("CREDHUB_CRED_PATH"), "CredHub Credential Path")
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
		ClientID:             cc.newClientConfig.clientIndentity,
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
			cc.credhubConfig.credPath, "uaa-client:"+cc.newClientConfig.clientIndentity,
			cc.credhubConfig.credPermissions,
		)

		if err != nil {
			cc.log.Error("Failed to add CredHub permission", err)
			return err
		}
	}
	return nil
}
