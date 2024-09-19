package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"onyx/types"
	"onyx/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install all packages listed in a package.json file",
	Run:   runDeps,
}

func runDeps(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Println("Error: package.json not found. Make sure you're in the correct directory.")
	}

	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		if err := os.Mkdir("node_modules", 0755); err != nil {
			fmt.Printf("Error creating node_modules directory: %v\n", err)
		}
	}

	fmt.Println("Reading package.json...")
	var packageJSON types.PackageJSON
	file, err := os.ReadFile("package.json")
	if err != nil {
		log.Fatalf("Error reading package.json: %v", err)
	}
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		log.Fatalf("Error parsing package.json: %v", err)
	}

	fmt.Println("🪄 Installing dependencies...")
	for name, version := range packageJSON.Dependencies {
		fmt.Printf("Installing %s@%s...\n", name, version)
		var versionNumber string

		if strings.HasPrefix(version, "^") {
			versionNumber = "latest"
		} else if strings.HasPrefix(version, "~") {
			versionNumber = version[1:]
		} else {
			versionNumber = version
		}

		tarballURL, resolvedVersion, err := util.GetPackageMetadata(name, versionNumber, remoteRegistry)
		if err != nil {
			log.Fatalf("Error fetching metadata for %s: %v", name, err)
		}

		err = util.DownloadAndExtract(tarballURL, name, "./node_modules")
		if err != nil {
			log.Fatalf("Error installing %s: %v", name, err)
		}

		err = util.UpdatePackageJSON(name, resolvedVersion, false)
		if err != nil {
			log.Fatalf("Error updating package.json: %v", err)
		}
	}

	fmt.Println("\n🛠️ Installing devDependencies...")
	for name, version := range packageJSON.DevDependencies {
		var versionNumber string

		if strings.HasPrefix(version, "^") {
			versionNumber = "latest"
		} else if strings.HasPrefix(version, "~") {
			versionNumber = version[1:]
		} else {
			versionNumber = version
		}

		tarballURL, resolvedVersion, err := util.GetPackageMetadata(name, versionNumber, remoteRegistry)
		if err != nil {
			log.Fatalf("Error fetching metadata for %s: %v", name, err)
		}

		err = util.DownloadAndExtract(tarballURL, name, "./node_modules")
		if err != nil {
			log.Fatalf("Error installing %s: %v", name, err)
		}

		err = util.UpdatePackageJSON(name, resolvedVersion, true)
		if err != nil {
			log.Fatalf("Error updating package.json: %v", err)
		}
	}

	fmt.Println("All packages installed successfully.")

}

func init() {
	rootCmd.AddCommand(depsCmd)
}
