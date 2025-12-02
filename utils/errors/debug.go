package errors

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func Debug() {
	// 1. VERIFY REGION AND PROFILE
	// Change these to match exactly what you use in your AWS CLI/Console
	targetRegion := "us-east-1"
	targetProfile := "staging"

	fmt.Printf("--- DIAGNOSTIC START ---\n")
	fmt.Printf("Target Profile: %s\n", targetProfile)
	fmt.Printf("Target Region:  %s\n", targetRegion)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(targetProfile),
		config.WithRegion(targetRegion),
	)
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	svc := ec2.NewFromConfig(cfg)

	// 2. NO FILTERS - GET EVERYTHING
	// We pass 'nil' filters to see absolutely every instance in this region
	input := &ec2.DescribeInstancesInput{}

	paginator := ec2.NewDescribeInstancesPaginator(svc, input)

	instanceCount := 0
	cutoffTime := time.Now().Add(-10 * time.Hour)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("API Error: %v", err)
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				instanceCount++

				// Calculate Age
				uptime := time.Since(*instance.LaunchTime)
				hours := int(uptime.Hours())

				// Check Age Condition
				isOlder := instance.LaunchTime.Before(cutoffTime)
				status := "[TOO NEW]"
				if isOlder {
					status = "[MATCHES TIME]"
				}

				fmt.Println("------------------------------------------------")
				fmt.Printf("ID: %s | State: %s | Region: %s\n", *instance.InstanceId, instance.State.Name, targetRegion)
				fmt.Printf("Launched: %s (%d hours ago) -> %s\n", instance.LaunchTime.Format("2006-01-02 15:04"), hours, status)

				// 3. PRINT ALL TAGS RAW
				// This is the most important part. See exactly what the SDK sees.
				fmt.Println("Tags found:")
				hasOpenCPETag := false
				for _, t := range instance.Tags {
					fmt.Printf("   Key: '%s' | Value: '%s'\n", *t.Key, *t.Value)
					if *t.Key == "OpenCPE" {
						hasOpenCPETag = true
					}
				}

				if !hasOpenCPETag {
					fmt.Println("   !! WARNING: Tag Key 'OpenCPE' NOT found on this instance !!")
				}
			}
		}
	}

	if instanceCount == 0 {
		fmt.Println("\n!!! ZERO INSTANCES FOUND !!!")
		fmt.Println("This means you are definitely querying the wrong REGION or ACCOUNT.")
		fmt.Println("Check your ~/.aws/config or the 'targetRegion' variable in this code.")
	}
}
