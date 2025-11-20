package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Flag declaration
var flagPolicy string
var flagRegion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opencpe",
	Short: "All-in-one Cloud-Native Policy Engine with sensible defaults",
	Long: `OpenCPE is a tool for managing cloud resources with heavily opinionated defaults.

Note: All Global flags are required
	
A comprehensive policy reference and usage instructions can be found at https://github.com/bazgab/opencpe
`,
	Run: func(cmd *cobra.Command, args []string) {
		// In case there are no flags being passed, return the help page
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		fmt.Println("Initializing: ")
		fmt.Printf("Policy: %s\n", flagPolicy)
		fmt.Printf("Region: %s\n", flagRegion)
		// You can call another function here, e.g., defaultRunLogic()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&flagRegion, "region", "", "Region to be queried (default is $OPENCPE_DEFAULT_REGION)")
	rootCmd.PersistentFlags().StringVar(&flagPolicy, "policy", "", "Policy to be executed")

	err := rootCmd.MarkPersistentFlagRequired("region")
	if err != nil {
		panic(err)
	}

	err = rootCmd.MarkPersistentFlagRequired("policy")
	if err != nil {
		panic(err)
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}
