package cmd

import (
	"fmt"
	"log"
	"onyx/types"
	"onyx/util"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const remoteRegistry = "https://registry.npmjs.org"

var getCmd = &cobra.Command{
	Use:   "get [package] [version]",
	Short: "Manually download and install a single npm package",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runGet(rootCmd, args, false)
	},
}

func runGet(cmd *cobra.Command, args []string, isTransitiveDependency bool) {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Println("Error: package.json not found. Make sure you're in the correct directory.")
		return
	}

	if _, err := os.Stat(".onyxlock.yaml"); os.IsNotExist(err) {
		err := util.CreateFile(".onyxlock.yaml")
		if err != nil {
			fmt.Println("Error: Couldn't create the lockfile")
			return
		}

		fmt.Println("🔐 Created .onyxlock.yaml")
	}

	packageName := args[0]
	version := "latest"
	if len(args) == 2 {
		version = args[1]
	}

	version = strings.TrimPrefix(version, "^")
	version = strings.TrimPrefix(version, "~")

	reGreater := regexp.MustCompile(`>=?\s*(\d+\.\d+\.\d+)`)
	greaterMatches := reGreater.FindStringSubmatch(version)

	if len(greaterMatches) > 1 {
		version = greaterMatches[1]
	}

	reLesser := regexp.MustCompile(`<=?\s*(\d+\.\d+\.\d+)`)
	lesserMatches := reLesser.FindStringSubmatch(version)

	if len(lesserMatches) > 1 {
		version = lesserMatches[1]
	}

	isDev, _ := cmd.Flags().GetBool("dev")
	isGlobal, _ := cmd.Flags().GetBool("global")
	if isGlobal {
		util.InstallGlobally(packageName, version, remoteRegistry)
		return
	}

	if !isTransitiveDependency {
		fmt.Printf("📦 Installing %s@%s...\n", packageName, version)
	}

	tarballURL, resolvedVersion, deps, err := util.GetPackageMetadata(packageName, version, remoteRegistry)
	if err != nil {
		log.Fatalf("Error fetching metadata for %s: %v", packageName, err)
	}

	err = util.DownloadAndExtract(tarballURL, packageName, "./node_modules")
	if err != nil {
		log.Fatalf("Error installing %s: %v", packageName, err)
	}

	err = util.UpdateLockfile(types.LockfileEntry{
		Version:  version,
		Resolved: resolvedVersion,
	}, packageName)

	if !isTransitiveDependency {
		err = util.UpdatePackageJSON(packageName, resolvedVersion, isDev)
		if err != nil {
			log.Fatalf("Error updating package.json: %v", err)
		}
	}

	for pkg, v := range deps {
		dependencyArgs := []string{pkg, fmt.Sprintf("%v", v)}
		runGet(cmd, dependencyArgs, true)
	}

	if !isTransitiveDependency {
		fmt.Printf("%s@%s installed successfully.\n", packageName, resolvedVersion)
	}

	if err != nil {
		fmt.Println("Warning: Couldn't update the .onyxlock.yaml file")
	}
}

func init() {
	getCmd.Flags().BoolP("dev", "d", false, "Install package as a devDependency")
	getCmd.Flags().BoolP("global", "g", false, "Install package globally")
	rootCmd.AddCommand(getCmd)
}
