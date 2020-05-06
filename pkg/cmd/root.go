package cmd

import (
	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces"
	"github.com/spf13/cobra"
)

type uaaApiFactory func(target string, zoneID string, adminClientIdentity string, adminClientSecret string) (interfaces.UaaAPI, error)

func uaaApiFactoryDefault(target string, zoneID string, adminClientIdentity string, adminClientSecret string) (interfaces.UaaAPI, error) {
	return interfaces.NewUaaApi(target, zoneID, adminClientIdentity, adminClientSecret)
}

type credHubFactory func(target string, skipTLS bool, clientID string, clientSecret string, uaaEndpoint string) (interfaces.CredHubAPI, error)

func credHubFactoryDefault(target string, skipTLS bool, clientID string, clientSecret string, uaaEndpoint string) (interfaces.CredHubAPI, error) {
	return interfaces.NewCredHubApi(target, skipTLS, clientID, clientSecret, uaaEndpoint)
}

func NewRootCmd(args []string) *cobra.Command {
	root := &cobra.Command{
		Use:     "uaaclient",
		Short:   "uaa-client-crud",
		Version: "0.1.0",
	}

	flags := root.PersistentFlags()
	out := root.OutOrStdout()
	root.AddCommand(
		NewCreateClientCmd(uaaApiFactoryDefault, credHubFactoryDefault, out),
		NewDeleteClientCmd(uaaApiFactoryDefault, credHubFactoryDefault, out),
	)

	flags.Parse(args)

	return root
}
