package cmd

import (
	"fmt"
	"onyx/util"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "Remove an already downloaded package",
	Args:  cobra.ExactArgs(1),
	Run:   runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {
	packageName := args[0]

	err := util.RemovePackageFromNodeModules(packageName)
	if err != nil {
		fmt.Printf("Error removing package from node_modules: %v", err)
	}

	err = util.RemovePackageFromPackageJSON(packageName)
	if err != nil {
		fmt.Printf("Error updating package.json: %v", err)
	}

	fmt.Printf("Package %s successfully removed.\n", packageName)

}

func init() {
	rootCmd.AddCommand(removeCmd)
}
