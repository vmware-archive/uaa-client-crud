package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(args []string) *cobra.Command {
	root := &cobra.Command{
		Use:     "uaaclient",
		Short:   "uaa-client-crud",
		Version: "0.1.0",
	}

	flags := root.PersistentFlags()
	out := root.OutOrStdout()
	root.AddCommand(
		NewCreateClientCmd(UaaApiFactoryDefault, out),
		NewDeleteClientCmd(out),
	)

	flags.Parse(args)

	return root
}
