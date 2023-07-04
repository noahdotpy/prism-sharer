package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/noahdotpy/prism-sharer/util"
)

// NOTE: Can we make this atomic?
// TODO: Check for existing symlinks to the store, remove if not declared
func ApplyRun(cmd *cobra.Command, args []string) {
	for groupName, group := range Config.Groups {
		for _, resource := range group.Resources {
			groupStore := filepath.Join(Config.StoresDir, groupName)
			resourcePath := filepath.Join(groupStore, "/", resource)

			if !util.DoesFileExist(resourcePath) {
				log.Warnf("Skipping to next resource as '%v' doesn't exist in the '%v' store.", resource, groupName)
				continue
			}

			for _, instanceName := range group.Instances {
				instanceDir := filepath.Join(Config.InstancesDir, "/", instanceName, "/.minecraft/")
				linkPath := filepath.Join(instanceDir, "/", resource)

				if !util.DoesFileExist(instanceDir) {
					log.Warnf("Skipping to next instance as '%v' doesn't exist.", instanceName)
					continue
				}

				linkPathInfo, err := os.Lstat(linkPath)
				blockerExists := err == nil

				if !blockerExists {
					createSymlink(resourcePath, linkPath, resource)
					continue
				}

				blockerIsSymlink := util.IsSymlink(linkPathInfo)

				var isLinkingBlocked bool

				if !blockerIsSymlink {
					isLinkingBlocked = true
				}

				if blockerIsSymlink {
					if isBlockingSymlinkExpected(resourcePath, linkPath) {
						log.Debugf("Skipped creating symlink for '%v' in instance '%v'. Expected symlink already exists.", resource, instanceName)
						continue
					}

					isLinkingBlocked = true
				}

				if isLinkingBlocked {
					blockerIsDir := linkPathInfo.IsDir()
					var blockerType string
					if blockerIsDir {
						blockerType = "folder"
					} else {
						blockerType = "file"
					}
					log.Warnf(
						"There is a %v blocking the creation of symlink for '%v' in instance '%v'.",
						blockerType,
						resource,
						instanceName)
				}
			}
		}
	}
}

func createSymlink(originPath string, linkPath string, resourceName string) error {

	log.Debugf("Creating link at '%v', with origin '%v'", linkPath, originPath)

	originInfo, err := os.Lstat(originPath)
	if err != nil {
		return err
	}

	var parentDirToCreate string
	if originInfo.IsDir() {
		parentDirToCreate = strings.TrimSuffix(linkPath, resourceName)
	} else {
		splittedPath := strings.SplitAfter(linkPath, "/")
		parentDirToCreate = strings.Join(splittedPath[:len(splittedPath)-1], "")
	}

	if !util.DoesFileExist(parentDirToCreate) {
		os.MkdirAll(parentDirToCreate, os.ModePerm)
	}

	err = os.Symlink(originPath, linkPath)
	if err != nil {
		log.Errorf("Error creating symlink: %v", err)
	}
	return nil
}

func isBlockingSymlinkExpected(expectedOrigin string, linkPath string) bool {
	actualOrigin, err := os.Readlink(linkPath)
	if err != nil {
		log.Fatalf("Could not read symlink: %v", err)
	}

	return expectedOrigin == actualOrigin

}
