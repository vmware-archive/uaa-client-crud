package main

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/lager"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
	"os"
)

type uaaClientConfig struct {
	uaaEndpoint         string
	adminClientIdentity string
	adminClientPwd      string
	clientIndentity     string
	clientPwd           string
}

type uaaClient struct {
	config uaaClientConfig
}

func newRootCmd(args []string) *cobra.Command {
	root := &cobra.Command{
		Use:   "uaaClientThing",
		Short: "Our Short String",
	}

	flags := root.PersistentFlags()
	//out := root.OutOrStdout()

	client := &uaaClient{}

	uaaCreateClientCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new client in UAA",
		RunE: func(cmd *cobra.Command, args []string) error {
			return client.uaaCreate()
		},
	}

	uaaCreateClientCmd.Flags().StringVarP(&client.config.uaaEndpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "uaa-endpoint")

	uaaCreateClientCmd.Flags().StringVarP(&client.config.adminClientIdentity, "admin-identity", "i", "", "Admin Username")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "admin-identity")

	uaaCreateClientCmd.Flags().StringVarP(&client.config.adminClientPwd, "admin-pwd", "p", "", "Admin Password")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "admin-pwd")

	uaaCreateClientCmd.Flags().StringVarP(&client.config.clientIndentity, "client-identity", "c", "", "New Client Identity")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "client-identity")

	uaaCreateClientCmd.Flags().StringVarP(&client.config.clientPwd, "client-pwd", "w", "", "New Client Password")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "client-pwd")

	root.AddCommand(uaaCreateClientCmd)

	flags.Parse(args)

	return root
}

func (cc *uaaClient) uaaCreate() error {

	fmt.Print("the endpoint is ")
	fmt.Println(cc.config.uaaEndpoint)
	logger := lager.NewLogger("test")

	// construct the API, and validate it
	api := uaa.New(cc.config.uaaEndpoint, "").WithClientCredentials(cc.config.adminClientIdentity, cc.config.adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)
	err := api.Validate()
	if err != nil {
		logger.Info(err.Error())
	}

	client := uaa.Client{
		ClientID:             cc.config.clientIndentity,
		ClientSecret:         cc.config.clientPwd,
		AccessTokenValidity:  1209600,
		AuthorizedGrantTypes: []string{"client_credentials", "refresh_token"},
		Scope:                []string{"openid", "oauth.approvals", "credhub.read", "credhub.write"},
		Authorities:          []string{"oauth.login", "credhub.read", "credhub.write"},
	}

	newClient, err := api.CreateClient(client)

	logger.Info(newClient.DisplayName)

	//api.ChangeClientSecret("credhub_admin_client", "iFB7oFXyRI1Yp3sHd_5RZ7WLDZHv2UX3")
	credHubClient, err := api.GetClient("credhub_admin_client")

	chAdmin, err := credhub.New("https://credhub-proxy.apps.brea.cf-app.com",
		credhub.SkipTLSValidation(true),
		credhub.Auth(auth.UaaClientCredentials("credhub_admin_client", credHubClient.ClientSecret)),
		credhub.AuthURL("https://uaa.sys.brea.cf-app.com"),
	)

	_, err = chAdmin.AddPermission(
		"/*", "uaa-client:"+cc.config.clientIndentity,
		[]string{"read", "write", "delete", "read_acl", "write_acl"},
	)

	return nil
}

func main() {
	fmt.Println("Hello")
	command := newRootCmd(os.Args[1:])
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
