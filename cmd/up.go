package cmd

import (
	"bufio"
	"fmt"
	"onyx/types"
	"onyx/util"
	"os"

	"github.com/spf13/cobra"
)

var createWithDefaultParameters bool

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Creates a new Node.js project",
	Long:  "This command initializes a new Node.js project by creating a package.json file.",
	Run:   runUp,
}

func runUp(cmd *cobra.Command, args []string) {
	pkg := types.PackageJSON{
		Version: "1.0.0",
		Name:    util.GetDefaultPackageName(),
		Main:    "index.js",
		License: "ISC",
	}

	if !createWithDefaultParameters {
		reader := bufio.NewReader(os.Stdin)

		pkg.Name = util.Prompt(reader, "package name: ", pkg.Name)
		pkg.Version = util.Prompt(reader, "version: ", pkg.Version)
		pkg.Description = util.Prompt(reader, "description: ", "")
		pkg.Main = util.Prompt(reader, "entry point: ", pkg.Main)
		pkg.Author = util.Prompt(reader, "author: ", "")
		pkg.License = util.Prompt(reader, "license: ", pkg.License)
	}

	err := util.WritePackageJSON(pkg)
	if err != nil {
		fmt.Printf("Error writing package.json: %v\n", err)
		return
	}

	fmt.Println("package.json file has been created.")
}

func init() {
	upCmd.Flags().BoolVarP(&createWithDefaultParameters, "skip", "y", false, "Skip prompts and use default values")
	rootCmd.AddCommand(upCmd)
}
