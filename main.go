package main

import (
	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	client "code.cloudfoundry.org/uaa-go-client"
	"code.cloudfoundry.org/uaa-go-client/config"
	"code.cloudfoundry.org/uaa-go-client/schema"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type uaaCreateCmd struct {
	uaaEndpoint         string
	adminClientIdentity string
	adminClientPwd      string
	clientIndentity     string
	clientPwd           string
}

func newRootCmd(args []string) *cobra.Command {
	root := &cobra.Command{
		Use:   "uaaClientThing",
		Short: "Our Short String",
	}

	flags := root.PersistentFlags()
	//out := root.OutOrStdout()

	cc := &uaaCreateCmd{

	}

	uaaCreateClientCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new client in UAA",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.uaaCreate()
		},
	}

	uaaCreateClientCmd.Flags().StringVarP(&cc.uaaEndpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "uaa-endpoint")

	uaaCreateClientCmd.Flags().StringVarP(&cc.adminClientIdentity, "admin-identity", "i", "", "Admin Username")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "admin-identity")

	uaaCreateClientCmd.Flags().StringVarP(&cc.adminClientPwd, "admin-pwd", "p", "", "Admin Password")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "admin-pwd")

	uaaCreateClientCmd.Flags().StringVarP(&cc.clientIndentity, "client-identity", "c", "", "New Client Identity")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "client-identity")

	uaaCreateClientCmd.Flags().StringVarP(&cc.clientPwd, "client-pwd", "w", "", "New Client Password")
	cobra.MarkFlagRequired(uaaCreateClientCmd.Flags(), "client-pwd")

	root.AddCommand(uaaCreateClientCmd)

	flags.Parse(args)

	return root
}

func (cc *uaaCreateCmd) uaaCreate() error {
	//
	fmt.Print("the endpoint is ")
	fmt.Println(cc.uaaEndpoint)

	cfg := &config.Config{
		ClientName:       cc.adminClientIdentity,
		ClientSecret:     cc.adminClientPwd,
		UaaEndpoint:      cc.uaaEndpoint,
		SkipVerification: true,
	}

	logger := lager.NewLogger("test")
	clock := clock.NewClock()

	uaaClient, err := client.NewClient(logger, cfg, clock)
	if  err != nil {
		return err
	}

	fmt.Printf("Connecting to: %s ...\n", cfg.UaaEndpoint)

	token, err := uaaClient.FetchToken(true)
	if err != nil{
		return err
	}

	fmt.Printf("Response:\n\ttoken: %s\n\texpires: %d\n", token.AccessToken, token.ExpiresIn)


	nc := &schema.OauthClient{
		ClientId:     cc.clientIndentity,
		ClientSecret: cc.clientPwd,
		Name: cc.clientIndentity,
		AuthorizedGrantTypes: []string{"refresh_token","client_credentials"},
		Authorities: []string{"cloud_controller.admin_read_only"},

	}
	newClient, err := uaaClient.RegisterOauthClient(nc)
	if err!= nil {
		return err
	}

	fmt.Printf("New client name is %s", newClient.Name)

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
