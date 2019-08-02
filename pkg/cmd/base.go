package cmd

import (
	"code.cloudfoundry.org/lager"
	"github.com/spf13/cobra"
	"io"
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

type uaaClientConfig struct {
	clientIndentity     string
	clientPwd           string
	clientGrantTypes    []string
	clientScopes        []string
	clientAuthorities   []string
	clientTokenValidity int64
}

type baseCmd struct {
	uaaConfig       uaaConfig
	credhubConfig   credhubConfig
	newClientConfig uaaClientConfig
	log             lager.Logger
	out             io.Writer
	verbose         bool
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
}

func (b *baseCmd) PreRun(cmd *cobra.Command, args []string) {
	if b.verbose {
		b.log.RegisterSink(lager.NewWriterSink(b.out, lager.DEBUG))
	} else {
		b.log.RegisterSink(lager.NewWriterSink(b.out, lager.ERROR))
	}
}
