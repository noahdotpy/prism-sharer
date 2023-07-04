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

// TODO: This whole function is spaghetti code, future me can fix it
// FIXME: It panics somewhere, fix it
func applyRun(cmd *cobra.Command, args []string) {
	// TODO: Check for existing symlinks to the store, removing if not declared

	for groupName, group := range loadedConfig.Groups {
		log.Debugf("in group: %v", groupName)
		for _, instanceName := range group.Instances {
			instanceDir := loadedConfig.InstancesDir + "/" + instanceName
			groupStoreDir := loadedConfig.StoreDir + groupName

			for _, toShare := range group.Shared {
				linkOrigin := filepath.Join(groupStoreDir, "/", toShare)
				createLinkAt := filepath.Join(instanceDir, "/", toShare)
				fileInfo, err := os.Lstat(createLinkAt)
				if err != nil {
					// No file
					log.Debugf("creating link at '%v', with origin '%v'", createLinkAt, linkOrigin)
					err = os.Symlink(linkOrigin, createLinkAt)
					if err != nil {
						log.Fatalf("Error creating symlink: %v", err)
					}
				}
				if err == nil {
					// File exists here

					if fileInfo.Mode()&os.ModeSymlink != 0 {
						// File is a symlink
						actualOrigin, err := os.Readlink(createLinkAt)
						if err != nil {
							log.Fatalf("Could not read symlink: %v", err)
						}

						if !(actualOrigin == linkOrigin) {
							log.Warnf("Could not create symlink for '%v' in instance '%v'. There is already a file here.", toShare, instanceName)
						}
					} else {
						log.Warnf("Could not create symlink for '%v' in instance '%v'. There is already a file here.", toShare, instanceName)
					}
				}
			}
		}
	}
}
