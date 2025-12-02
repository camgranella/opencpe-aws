package errors

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"
	"os"
)

func IdentityCheck(profile string, region string, account string) {
	targetProfile := profile
	targetRegion := region
	expectedAccount := account

	fmt.Println("------ IDENTITY CHECK ------")
	fmt.Println("")
	fmt.Println("------ ENVIRONMENT DIAGNOSTICS ------")
	if os.Getenv("AWS_PROFILE") == "" {
		fmt.Println("AWS_PROFILE environment variable not set.")
	} else {
		fmt.Printf("AWS_PROFILE env var: '%s'\n", os.Getenv("AWS_PROFILE"))
	}
	if os.Getenv("AWS_REGION") == "" {
		fmt.Println("AWS_REGION environment variable not set.")
	} else {
		fmt.Printf("AWS_REGION env var: '%s'\n", os.Getenv("AWS_REGION"))
	}
	fmt.Println("-------------------------------------")

	// Load Config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(targetProfile),
		config.WithRegion(targetRegion),
	)
	if err != nil {
		log.Fatalf("Config Load Error: %v", err)
	}

	client := sts.NewFromConfig(cfg)

	// Ask AWS: "Who am I?"
	identity, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("Authentication Failed: %v\n(Check if you ran 'aws sso login --profile staging')", err)
	}

	currentAccount := *identity.Account

	fmt.Println("------ AUTHENTICATION RESULT ------")
	fmt.Printf("Connected Account ID: %s\n", currentAccount)
	fmt.Printf("Expected Account ID:  %s\n", expectedAccount)
	fmt.Printf("Connected ARN:        %s\n", *identity.Arn)
	fmt.Println("-----------------------------------")

	if currentAccount != expectedAccount {
		fmt.Println("❌ MISMATCH DETECTED!")
		fmt.Println("You are NOT connected to the account in your config file.")
		fmt.Println("Your terminal is likely pointing to a different set of credentials.")
	} else {
		fmt.Println("✅ Account ID matches.")
		fmt.Println("If this matches, but you see 0 instances, verify the INSTANCE ID in the console.")
		fmt.Println("Does the Instance ID exist in region us-east-1 of THIS specific account?")
	}
}
