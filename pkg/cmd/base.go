package cmd

import (
	"io"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/spf13/cobra"
)

type uaaConfig struct {
	endpoint            string
	adminClientIdentity string
	adminClientPwd      string
}

type credhubConfig struct {
	endpoint        string
	clientID        string
	clientPwd       string
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
}

func newBaseCmd(out io.Writer) baseCmd {
	base := baseCmd{
		out: out,
		log: lager.NewLogger("uaa-crud-cli"),
	}

	return base
}

func (b *baseCmd) addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.verbose, "verbose", "v", false, "verbose logging")
	cmd.Flags().StringVarP(&b.uaaConfig.endpoint, "uaa-endpoint", "e", "", "UAA Endpoint")
	cmd.MarkFlagRequired("uaa-endpoint")
	cmd.Flags().StringVarP(&b.uaaConfig.adminClientIdentity, "admin-identity", "u", "", "UAA Admin Username")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-identity")
	cmd.Flags().StringVarP(&b.uaaConfig.adminClientPwd, "admin-pwd", "p", "", "UAA Admin Password")
	cobra.MarkFlagRequired(cmd.Flags(), "admin-pwd")
	cmd.Flags().StringVarP(&b.targetClientIdentity, "target-client-identity", "c", "", "Target Client Identity")
	cobra.MarkFlagRequired(cmd.Flags(), "target-client-identity")

	cmd.Flags().StringVar(&b.credhubConfig.clientID, "credhub-identity", os.Getenv("CREDHUB_CLIENT_ID"), "Credhub Client Identity if granting or revoking CredHub access")
	cmd.Flags().StringVar(&b.credhubConfig.clientPwd, "credhub-password", os.Getenv("CREDHUB_CLIENT_PASSWORD"), "Credhub Client Password if granting or revoking CredHub access")
	cmd.Flags().StringVar(&b.credhubConfig.endpoint, "credhub-endpoint", os.Getenv("CREDHUB_URL"), "CredHub endpoint URL")
	cmd.Flags().StringVar(&b.credhubConfig.credPath, "credential-path", os.Getenv("CREDHUB_CRED_PATH"), "CredHub Credential Path")

}

func (b *baseCmd) PreRun(cmd *cobra.Command, args []string) {
	if b.verbose {
		b.log.RegisterSink(lager.NewWriterSink(b.out, lager.DEBUG))
	} else {
		b.log.RegisterSink(lager.NewWriterSink(b.out, lager.ERROR))
	}
}
