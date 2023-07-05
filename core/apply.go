package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/noahdotpy/prism-sharer/util"
)

// NOTE: Can we make this atomic?
func ApplyRun(cmd *cobra.Command, args []string) {
	for groupName, group := range Config.Groups {
		groupStore := filepath.Join(Config.StoresDir, groupName)
		for _, instanceName := range group.Instances {
			instanceDir := filepath.Join(Config.InstancesDir, "/", instanceName, "/.minecraft/")

			log.Debug("Starting pre-applying cleanup", "instance", instanceName)

			err := filepath.WalkDir(instanceDir, func(path string, d fs.DirEntry, err error) error {
				fileInfo, err := os.Lstat(path)
				if err != nil {
					return err
				}

				if !util.IsSymlink(fileInfo) {
					return nil
				}
				symlinkOrigin, err := os.Readlink(path)
				if err != nil {
					return err
				}
				isLinkDeclared := slices.ContainsFunc(group.Resources, func(el string) bool {
					return strings.HasPrefix(symlinkOrigin, groupStore)
				})
				if !isLinkDeclared {
					log.Warn("Removing old link", "path", path)
					os.Remove(path)
				}
				// TODO: Is it possible to remove leftover parent directories of deleted links
				log.Debugf("'%v' symlink is declared: %v", symlinkOrigin, isLinkDeclared)
				return nil
			})
			if err != nil {
				log.Errorf("impossible to walk directories: %s", err)
			}
		}
		for _, resource := range group.Resources {
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

				log.Debug("There is a file blocking symlink creation", "resource", resource, "instanceName", instanceName)

				blockerIsSymlink := util.IsSymlink(linkPathInfo)
				log.Debug("", "    blockerIsSymlink", blockerIsSymlink)

				var isLinkingBlocked bool

				if !blockerIsSymlink {
					isLinkingBlocked = true
				}

				if blockerIsSymlink {

					blockerOrigin, err := os.Readlink(linkPath)
					if err != nil {
						log.Errorf("Could not read symlink for '%v' in instance '%v': %v", resource, instanceName, err)
					}
					log.Debug("", "    blockerOrigin", blockerOrigin)
					if isBlockingSymlinkExpected(resourcePath, linkPath) {
						log.Debugf("     Skipped creating symlink, blocker has the expected origin")
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
	blockerOrigin, err := os.Readlink(linkPath)
	if err != nil {
		log.Fatalf("Could not read symlink: %v", err)
	}

	return expectedOrigin == blockerOrigin

}
