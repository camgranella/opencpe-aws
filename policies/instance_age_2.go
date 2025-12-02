package policies

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func InstanceAge2Days(profile string, region string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	svc := ec2.NewFromConfig(cfg)

	// 1. Set Cutoff: 10 Hours ago
	cutoffTime := time.Now().Add(-2 * time.Hour)

	fmt.Printf("Searching for instances launched BEFORE: %s\n", cutoffTime.Format(time.RFC3339))

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{"OpenCPE"}, // Ensure your instances actually have this tag!
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	}

	paginator := ec2.NewDescribeInstancesPaginator(svc, input)

	matchCount := 0

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page, %v", err)
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {

				// 2. Logic Check: Is LaunchTime BEFORE the cutoff?
				// "Older" than 10 hours means LaunchTime < (Now - 10h)
				if instance.LaunchTime != nil && instance.LaunchTime.Before(cutoffTime) {
					matchCount++

					// Calc hours for display
					uptime := time.Since(*instance.LaunchTime)
					hours := int(uptime.Hours())

					name := "N/A"
					for _, t := range instance.Tags {
						if *t.Key == "Name" {
							name = *t.Value
							break
						}
					}

					fmt.Printf("-- FOUND: %s | ID: %s | Uptime: %d hours\n",
						name, *instance.InstanceId, hours)
				}
			}
		}
	}

	if matchCount == 0 {
		fmt.Println("No instances found that are older than 10 hours with tag 'OpenCPE'.")
	}
}
