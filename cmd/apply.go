package cmd

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var (
	currOriginIsExpected bool

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply configuration changes",
		Run:   applyRun,
	}
)

func applyRun(cmd *cobra.Command, args []string) {
	// TODO: Check for existing symlinks to the store, removing if not declared

	for groupName, group := range loadedConfig.Groups {

		for _, instanceName := range group.Instances {
			instanceDir := filepath.Join(loadedConfig.InstancesDir, "/", instanceName)
			groupStoreDir := filepath.Join(loadedConfig.StoreDir, groupName)

			_, err := os.Lstat(instanceDir)
			instanceDirExists := err == nil
			if !instanceDirExists {
				log.Warnf("Instance '%v' doesn't exist, is this a typo?", instanceName)
				continue
			}

			for _, toShare := range group.Shared {

				linkOrigin := filepath.Join(groupStoreDir, "/", toShare)
				createLinkAt := filepath.Join(instanceDir, "/", toShare)

				log.Debug(linkOrigin)
				_, err := os.Lstat(linkOrigin)
				sharedResourceExists := err == nil
				if !sharedResourceExists {
					log.Warnf("'%v' doesn't exist in the '%v' store, consider making it.", toShare, groupName)
					continue
				}

				resourceInfo, err := os.Lstat(createLinkAt)

				resourceExists := err == nil

				if !resourceExists {

					log.Debugf("Creating link at '%v', with origin '%v'", createLinkAt, linkOrigin)
					err = os.Symlink(linkOrigin, createLinkAt)
					if err != nil {
						log.Warnf("Error creating symlink: %v", err)
					}
					continue
				}

				fileIsSymlink := resourceInfo.Mode()&os.ModeSymlink != 0

				if !fileIsSymlink {
					log.Warnf("There is a file blocking the creation of symlink for '%v' in instance '%v'.", toShare, instanceName)
					continue
				}

				expectedOrigin := linkOrigin
				actualOrigin, err := os.Readlink(createLinkAt)
				if err != nil {
					log.Fatalf("Could not read symlink: %v", err)
				}

				if actualOrigin != expectedOrigin {
					log.Warnf("Could not create symlink for '%v' in instance '%v'. There is already a file here.", toShare, instanceName)
				} else {
					log.Debugf("Skipped creating symlink for '%v' in instance '%v'. Expected symlink already exists.", toShare, instanceName)
				}
			}
		}
	}
}
