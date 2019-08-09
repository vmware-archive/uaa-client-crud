package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

type clientDeleteCmd struct {
	baseCmd
}

func NewDeleteClientCmd(uaaApiFactory uaaApiFactory, credHubFactory credHubFactory, out io.Writer) *cobra.Command {
	cd := &clientDeleteCmd{
		baseCmd: newBaseCmd(uaaApiFactory, credHubFactory, out),
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
	apiClient := cd.uaaApiFactory(cd.uaaConfig.endpoint, "", cd.uaaConfig.adminClientIdentity, cd.uaaConfig.adminClientPwd)
	err := apiClient.Validate()
	if err != nil {
		cd.log.Error("", err)
	}
	_, err = apiClient.GetClient(cd.targetClientIdentity)
	if err == nil {
		_, err = apiClient.DeleteClient(cd.targetClientIdentity)
		if err != nil {
			cd.log.Error("Failed to delete UAA client ["+cd.targetClientIdentity+"]", err)
			return err
		} else {
			cd.log.Debug("UAA client [" + cd.targetClientIdentity + "] deleted")
		}
	} else {
		cd.log.Debug("UAA client [" + cd.targetClientIdentity + "]. Skipping delete")
	}

	if cd.credhubConfig.endpoint != "" && cd.credhubConfig.clientID != "" && cd.credhubConfig.credPath != "" && cd.credhubConfig.clientPwd != "" {

		chAdmin, err := cd.credHubFactory(cd.credhubConfig.endpoint,
			true,
			cd.credhubConfig.clientID,
			cd.credhubConfig.clientPwd,
			cd.uaaConfig.endpoint,
		)

		if err != nil {
			cd.log.Error("Failed to connect to CredHub", err)
			return err
		}

		permission, err := chAdmin.GetPermissionByPathActor(cd.credhubConfig.credPath, "uaa-client:"+cd.targetClientIdentity)
		if err != nil {
			cd.log.Debug("Failed to get permission object from CredHub. Skipping delete")
		} else {
			_, err = chAdmin.DeletePermission(permission.UUID)
			if err != nil {
				cd.log.Error("Failed to delete permission object in CredHub", err)
				return err
			} else {
				cd.log.Debug("CredHub permission deleted")
			}
		}

	}
	return nil
}
