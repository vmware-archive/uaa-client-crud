package cmd

import (
	"io"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

type clientDeleteCmd struct {
	baseCmd
}

func NewDeleteClientCmd(out io.Writer) *cobra.Command {
	cd := &clientDeleteCmd{
		baseCmd: newBaseCmd(out),
	}

	cmd := &cobra.Command{
		Use:    "delete",
		Short:  "Delete a client in UAA",
		PreRun: cd.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return cd.run()
		},
	}

	cd.addCommonFlags(cmd)

	return cmd
}

func (cd *clientDeleteCmd) run() error {

	// construct the API, and validate it
	api := uaa.New(cd.uaaConfig.endpoint, "").WithClientCredentials(cd.uaaConfig.adminClientIdentity, cd.uaaConfig.adminClientPwd, uaa.JSONWebToken).WithSkipSSLValidation(true)
	err := api.Validate()
	if err != nil {
		cd.log.Error("", err)
	}

	_, err = api.DeleteClient(cd.targetClientIdentity)

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

		permission, err := chAdmin.GetPermissionByPathActor(cd.credhubConfig.credPath, "uaa-client:"+cd.targetClientIdentity)
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
