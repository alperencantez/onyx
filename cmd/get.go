package cmd

import (
	"fmt"
	"log"
	"onyx/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const remoteRegistry = "https://registry.npmjs.org"

var getCmd = &cobra.Command{
	Use:   "get [package] [version]",
	Short: "Manually download and install a single npm package",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runGet(rootCmd, args)
	},
}

func runGet(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Println("Error: package.json not found. Make sure you're in the correct directory.")
		return
	}

	packageName := args[0]
	version := "latest"
	if len(args) == 2 {
		version = args[1]
	}

	version = strings.TrimPrefix(version, "^")
	isDev, _ := cmd.Flags().GetBool("dev")
	isGlobal, _ := cmd.Flags().GetBool("global")
	if isGlobal {
		util.InstallGlobally(packageName, version, remoteRegistry)
		return
	}

	fmt.Printf("📦 Installing %s@%s...\n", packageName, version)
	tarballURL, resolvedVersion, deps, err := util.GetPackageMetadata(packageName, version, remoteRegistry)
	if err != nil {
		log.Fatalf("Error fetching metadata for %s: %v", packageName, err)
	}

	err = util.DownloadAndExtract(tarballURL, packageName, "./node_modules")
	if err != nil {
		log.Fatalf("Error installing %s: %v", packageName, err)
	}

	err = util.UpdatePackageJSON(packageName, resolvedVersion, isDev)
	if err != nil {
		log.Fatalf("Error updating package.json: %v", err)
	}

	for pkg, v := range deps {
		dependencyArgs := []string{pkg, fmt.Sprintf("%v", v)}
		runGet(cmd, dependencyArgs)
	}

	fmt.Printf("%s@%s installed successfully.\n", packageName, resolvedVersion)
}

func init() {
	getCmd.Flags().BoolP("dev", "d", false, "Install package as a devDependency")
	getCmd.Flags().BoolP("global", "g", false, "Install package globally")
	rootCmd.AddCommand(getCmd)
}
