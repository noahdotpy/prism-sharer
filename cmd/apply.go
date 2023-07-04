package cmd

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/noahdotpy/prism-sharer/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration changes",
	Run:   applyRun,
}

func applyRun(cmd *cobra.Command, args []string) {
	// TODO: Check for existing symlinks, removing if not declared

	for groupName, group := range loadedConfig.Groups {
		var group config.Group = group

		for _, instanceName := range group.Instances {
			instanceDir := loadedConfig.InstancesDir + "/" + instanceName
			groupStoreDir := loadedConfig.InstancesDir + groupName

			for _, toShare := range group.Shared {
				err := os.Symlink(groupStoreDir, filepath.Join(instanceDir+"/"+toShare))
				if err != nil {
					// TODO: Check if symlink resolves to what is wanted, if it does then don't print warning
					log.Warnf("Could not create symlink for '%v'. Did you move the existing file/folder?", toShare)
					log.Debug(err)
				} else {
					log.Infof("Succesfully created symlink for '%v'.", toShare)
				}
			}
		}
	}
}
