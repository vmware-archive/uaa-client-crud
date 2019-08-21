package cmd

import (
	"io"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/spf13/cobra"
)

type uaaConfig struct {
	endpoint            string
	adminClientIdentity string
	adminClientSecret   string
}

type credhubConfig struct {
	endpoint        string
	clientID        string
	clientSecret    string
	credPath        string
	credPermissions []string
}

type baseCmd struct {
	uaaConfig            uaaConfig
	credhubConfig        credhubConfig
	log                  lager.Logger
	out                  io.Writer
	verbose              bool
	targetClientIdentity string
	uaaApiFactory        uaaApiFactory
	credHubFactory       credHubFactory
}

func newBaseCmd(uaaApiFactory uaaApiFactory, credHubFactory credHubFactory, out io.Writer) baseCmd {
	base := baseCmd{
		out:            out,
		log:            lager.NewLogger("uaa-crud"),
		uaaApiFactory:  uaaApiFactory,
		credHubFactory: credHubFactory,
	}

	return base
}

func (b *baseCmd) addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.verbose, "verbose", "v", false, "verbose logging")
	cmd.Flags().StringVarP(&b.uaaConfig.endpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cmd.MarkFlagRequired("uaa-endpoint")
	cmd.Flags().StringVarP(&b.uaaConfig.adminClientIdentity, "admin-identity", "u", "", "UAA Admin Username")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-identity")
	cmd.Flags().StringVarP(&b.uaaConfig.adminClientSecret, "admin-secret", "p", "", "UAA Admin Client Secret")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-secret")
	cmd.Flags().StringVarP(&b.targetClientIdentity, "target-client-identity", "c", "", "Target Client Identity")
	cobra.MarkFlagRequired(cmd.Flags(), "target-client-identity")

	cmd.Flags().StringVar(&b.credhubConfig.clientID, "credhub-identity", os.Getenv("CREDHUB_CLIENT_ID"), "Credhub Client Identity if granting or revoking CredHub access")
	cmd.Flags().StringVar(&b.credhubConfig.clientSecret, "credhub-secret", os.Getenv("CREDHUB_CLIENT_SECRET"), "Credhub Client Password if granting or revoking CredHub access")
	cmd.Flags().StringVar(&b.credhubConfig.endpoint, "credhub-endpoint", os.Getenv("CREDHUB_URL"), "CredHub endpoint URL")
	cmd.Flags().StringVar(&b.credhubConfig.credPath, "credential-path", os.Getenv("CREDHUB_CRED_PATH"), "CredHub Credential Path")

}

func (b *baseCmd) PreRun(cmd *cobra.Command, args []string) {
	if b.verbose {
		b.log.RegisterSink(lager.NewWriterSink(b.out, lager.DEBUG))
	} else {
		b.log.RegisterSink(lager.NewWriterSink(b.out, lager.ERROR))
	}
	b.uaaConfig.endpoint = prependHttps(b.uaaConfig.endpoint)
	b.credhubConfig.endpoint = prependHttps(b.credhubConfig.endpoint)
}

func prependHttps(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http", "https", 1)
	} else if strings.HasPrefix(url, "https://") {
		return url
	} else if url == "" {
		return url
	} else {
		return "https://" + url
	}
}
