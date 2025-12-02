package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bazgab/opencpe/utils/logging"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type Config struct {
	Authentication struct {
		AwsProfile string `json:"aws_profile"`
	} `json:"authentication"`

	Notification struct {
		SmtpHost      string `json:"smtp_host"`
		SmtpPort      int    `json:"smtp_port"`
		EmailFrom     string `json:"email_from"`
		EmailPassword string `json:"email_password"`
	} `json:"notification"`

	IgnoredTags struct {
		Owner   []string `json:"owner"`
		Project []string `json:"project"`
	} `json:"ignored_tags"`
}

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "notify resource owners of policy infringement",
	Long:  `This command only notifies resource owners of policy infringement, as opposed to notify-and-delete.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if the global flag "policy" was actually set by the user
		RequiredFlags := []string{"config", "policy", "region"}
		for _, flag := range RequiredFlags {
			if !cmd.Flags().Changed(flag) {
				return fmt.Errorf("required flag %s not set for the 'notify' command", flag)
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Checking if global flags are working
		fmt.Println("Query Request Output:")
		logging.TextRequestOutputLogger("Printing values from global flags", flagConfig, flagPolicy, flagRegion)

		cfgFile, err := os.ReadFile(flagConfig)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}

		var config Config

		err = json.Unmarshal(cfgFile, &config)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}

		fmt.Println("Loaded Configuration:")
		fmt.Printf("-- Authentication.Profile: %s\n", config.Authentication.AwsProfile)
		fmt.Printf("-- Notification.SMTP Host: %s\n", config.Notification.SmtpHost)
		fmt.Printf("-- Notification.SMTP Port: %d\n", config.Notification.SmtpPort)
		fmt.Printf("-- Notification.From Email: %s\n", config.Notification.EmailFrom)
		for _, owner := range config.IgnoredTags.Owner {
			fmt.Printf("-- IgnoredTags.Owner: %s\n", owner)
		}
		for _, project := range config.IgnoredTags.Project {
			fmt.Printf("-- IgnoredTags.Project: %s\n", project)
		}

	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)

}
