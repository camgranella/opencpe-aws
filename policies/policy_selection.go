package policies

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type PolicyStruct struct {
	Policy     string
	AwsProfile string
	Region     string
	Config     string
}

type InstanceResult struct {
	Id     string
	Name   string
	Owner  string
	Uptime int
}

type AppConfig struct {
	IgnoredTags map[string][]string `json:"ignored_tags"`
}

func LoadConfig(filename string) (*AppConfig, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func sliceContains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
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

	case "instance-age-7-days":
		instances, err := InstanceAge7Days(p.AwsProfile, p.Region, appConfig.IgnoredTags)
		if err != nil {
			log.Fatal(err)
		}
		return instances

	}
	return nil
}
