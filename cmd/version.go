package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	Version   = "0.0.1 - development"
	BuildDate = "2025-12-1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Extended OpenCPE version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OpenCPE Version: %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Go Version: %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
