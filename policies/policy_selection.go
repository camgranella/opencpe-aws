package policies

import (
	"fmt"
	"log"
)

type PolicyStruct struct {
	Policy     string
	AwsProfile string
	Region     string
	Config     string
}

func SelectPolicy(p PolicyStruct) []InstanceResult {

	fmt.Println("")
	fmt.Printf("Policy: %s\n", p.Policy)
	fmt.Printf("Profile: %s\n", p.AwsProfile)
	fmt.Printf("Region: %s\n", p.Region)

	appConfig, _ := LoadConfig(p.Config)

	switch p.Policy {
	case "instance-age-2-days":
		instances, err := InstanceAge2Days(p.AwsProfile, p.Region, appConfig.IgnoredTags)
		if err != nil {
			log.Fatal(err)
		}
		return instances

	case "other-policy":
		instances, err := InstanceAge2Days(p.AwsProfile, p.Region, appConfig.IgnoredTags)
		if err != nil {
			log.Fatal(err)
		}
		return instances

	}
	return nil
}
