package policies

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"strings"
	"time"
)

// Update signature to accept the ignoredTags map
func InstanceAge7Days(profile string, region string, ignoredTags map[string][]string) ([]InstanceResult, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := ec2.NewFromConfig(cfg)

	cutoffTime := time.Now().Add(-168 * time.Hour)

	fmt.Printf("Searching for instances launched BEFORE: %s\n", cutoffTime.Format(time.RFC3339))

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{Name: aws.String("tag-key"), Values: []string{"OpenCPE"}},
			{Name: aws.String("tag-key"), Values: []string{"Owner"}},
			{Name: aws.String("instance-state-name"), Values: []string{"running"}},
		},
	}

	var results []InstanceResult

	paginator := ec2.NewDescribeInstancesPaginator(svc, input)

	matchCount := 0

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to get page: %w", err)
		}

		fmt.Println("-- Found the following instances: ")

		for _, reservation := range page.Reservations {
		InstanceLoop:
			for _, instance := range reservation.Instances {

				if instance.LaunchTime != nil && instance.LaunchTime.Before(cutoffTime) {

					name := "N/A"
					owner := "N/A"
					for _, t := range instance.Tags {

						switch *t.Key {
						case "Name":
							name = *t.Value
						case "Owner":
							owner = *t.Value
						}

						for configKey, configValues := range ignoredTags {

							// Check if the AWS Tag Key matches the Config Key
							if strings.EqualFold(*t.Key, configKey) {

								// Check if the Value exists in the list
								if sliceContains(configValues, *t.Value) {
									// If found skip this instance
									continue InstanceLoop
								}
							}
						}
					}

					matchCount++
					uptime := time.Since(*instance.LaunchTime)
					hours := int(uptime.Hours())

					fmt.Printf("-- Instance Name: %s | ID: %s | Uptime: %d hours | Owner: %s\n",
						name, *instance.InstanceId, hours, owner)

					data := InstanceResult{
						Id:     *instance.InstanceId,
						Name:   name,
						Owner:  owner,
						Uptime: hours,
					}

					results = append(results, data)
				}
			}
		}
	}

	if matchCount == 0 {
		fmt.Println("No instances found that match criteria.")
	}

	return results, nil
}
