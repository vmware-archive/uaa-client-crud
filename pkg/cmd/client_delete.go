package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type clientDeleteCmd struct {
	baseCmd
	targetClientSecret string
	deleteCredhubPath  bool
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
	deleteCredhubPath, _ := strconv.ParseBool(os.Getenv("DELETE_CREDHUB_PATH"))
	cmd.Flags().BoolVar(&cd.deleteCredhubPath, "delete-credhub-path", deleteCredhubPath, "Delete all credentials in path")
	cmd.Flags().StringVarP(&cd.targetClientSecret, "target-client-secret", "w", "", "Target Client Secret")

	return cmd
}

func (cd *clientDeleteCmd) run() error {
	if cd.credhubConfig.endpoint != "" && cd.credhubConfig.clientID != "" && cd.credhubConfig.credPath != "" && cd.credhubConfig.clientSecret != "" {

		if cd.deleteCredhubPath {
			chClient, err := cd.credHubFactory(cd.credhubConfig.endpoint,
				true,
				cd.targetClientIdentity,
				cd.targetClientSecret,
				cd.uaaConfig.endpoint,
			)
			if err != nil {
				cd.log.Error(fmt.Sprintf("Failed to connect to credhub client for deleting credentials by path [%s]", cd.targetClientIdentity), err)
				return err
			}

			cd.log.Debug(fmt.Sprintf("Ch Client [%+v]", chClient))

			deletePath := strings.ReplaceAll(cd.credhubConfig.credPath, "*", "")
			credPaths, err := chClient.FindByPath(deletePath)
			if err != nil {
				cd.log.Error(fmt.Sprintf("Failed to lookup credentials by path [%s]", deletePath), err)
				return err
			}

			for _, credPath := range credPaths {
				err := chClient.DeleteCredential(credPath)
				if err != nil {
					cd.log.Error(fmt.Sprintf("Failed to delete credential [%s]", credPath), err)
					return err
				}
			}
		}

		chAdmin, err := cd.credHubFactory(cd.credhubConfig.endpoint,
			true,
			cd.credhubConfig.clientID,
			cd.credhubConfig.clientSecret,
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

	apiClient, err := cd.uaaApiFactory(cd.uaaConfig.endpoint, "", cd.uaaConfig.adminClientIdentity, cd.uaaConfig.adminClientSecret)
	if err != nil {
		cd.log.Error("Error validation UAA API client", err)
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
		cd.log.Error(fmt.Sprintf("Unable to Get UAA client [%s]", cd.targetClientIdentity), err)
		cd.log.Debug("UAA client [" + cd.targetClientIdentity + "]. Skipping delete")
	}

	return nil
}
