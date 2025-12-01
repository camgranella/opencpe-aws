package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "only notify resource owners",
	Long:  `This command only notifies resource owners of policy infringement, as opposed to notify-and-delete.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if the global flag "policy" was actually set by the user
		if !cmd.Flags().Changed("policy") {
			return fmt.Errorf("required flag 'policy' not set for the 'notify' command")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Checking if global flags are working
		fmt.Println("Priting value from global flags")
		fmt.Printf("Policy: %s\n", flagPolicy)
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)

}
