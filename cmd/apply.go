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

// NOTE: Can we make this atomic?
// TODO: Check for existing symlinks to the store, remove if not declared
func applyRun(cmd *cobra.Command, args []string) {
	for groupName, group := range loadedConfig.Groups {
		for _, resource := range group.Resources {
			groupStore := filepath.Join(loadedConfig.StoreDir, groupName)
			resourcePath := filepath.Join(groupStore, "/", resource)

			log.Debug(resourcePath)
			_, err := os.Lstat(resourcePath)
			resourceExists := err == nil
			if !resourceExists {
				log.Errorf("Resource '%v' doesn't exist in the '%v' store.", resource, groupName)
				continue
			}

			for _, instanceName := range group.Instances {
				instanceDir := filepath.Join(loadedConfig.InstancesDir, "/", instanceName, "/.minecraft/")
				createLinkAt := filepath.Join(instanceDir, "/", resource)

				_, err := os.Lstat(instanceDir)
				instanceDirExists := err == nil
				if !instanceDirExists {
					log.Errorf("Instance '%v' doesn't exist.", instanceName)
					continue
				}
				resourceInfo, err := os.Lstat(createLinkAt)

				resourceExists := err == nil

				if !resourceExists {

					log.Debugf("Creating link at '%v', with origin '%v'", createLinkAt, resourcePath)

					// TODO: Make directory if parent directory of resource doesn't exist
					// os.MkdirAll(DIR, os.ModePerm)

					err = os.Symlink(resourcePath, createLinkAt)
					if err != nil {
						log.Errorf("Error creating symlink: %v", err)
					}
					continue
				}

				fileIsSymlink := resourceInfo.Mode()&os.ModeSymlink != 0

				var resourceBlocked bool
				var expectedOrigin string
				var actualOrigin string

				if !fileIsSymlink {
					resourceBlocked = true
				}

				if fileIsSymlink {
					expectedOrigin = resourcePath
					actualOrigin, err = os.Readlink(createLinkAt)
					if err != nil {
						log.Fatalf("Could not read symlink: %v", err)
					}

					if expectedOrigin == actualOrigin {
						log.Debugf("Skipped creating symlink for '%v' in instance '%v'. Expected symlink already exists.", resource, instanceName)
						continue
					}

					resourceBlocked = true
				}

				if resourceBlocked {
					log.Warnf("There is a file blocking the creation of symlink for '%v' in instance '%v'.", resource, instanceName)
				}
			}
		}
	}
}
