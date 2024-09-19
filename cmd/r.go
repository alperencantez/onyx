package cmd

import (
	"fmt"
	"onyx/util"
	"os"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "r [script]",
	Short: "Runs a custom npm script defined in package.json",
	Args:  cobra.ExactArgs(1),
	Run:   runR,
}

func runR(cmd *cobra.Command, args []string) {
	scriptName := args[0]

	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Println("package.json not found")
	}

	packageJSON, err := util.ReadPackageJSON()
	if err != nil {
		fmt.Printf("Error reading package.json: %v\n", err)
		return
	}

	script, exists := packageJSON.Scripts[scriptName]
	if !exists {
		fmt.Printf("Script '%s' not found in package.json\n", scriptName)
		return
	}

	err = util.RunCustomScript(script)
	if err != nil {
		fmt.Printf("Error running script '%s': %v\n", scriptName, err)
		return
	}

	fmt.Printf("\n✨ Successfully ran script '%s'.\n", scriptName)
}

func init() {
	rootCmd.AddCommand(runCmd)
}
