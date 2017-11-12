package main

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jasonriedel/extractor/lib"
	"github.com/jasonriedel/extractor/sources/aws-cloud/services"
)


func TestGetRegions(t *testing.T) {
	regions := getRegions()
	if !(len(regions) > 0) {
		t.Error("No regions were found.")
	}
}

func TestSetupAwsSessions(t *testing.T) {
	accountsMap := make(map[string]string)
	accountsMap["tuxlabs"] = "907391580367"
	regions := []string{"us-east-1"}
	awsSessions := setupAwsSessions(accountsMap, regions)

	//not sure if this will ever return nil
	if awsSessions == nil {
		t.Error("Unable to create an AWS Session.")
	}
}

func TestCollectLightSailInstances(t *testing.T) {

	configuration, err := lib.LoadConfiguration(*fConfig)
	if err != nil {
		t.Error("Unable to load the configuration file!")
	}

	regions := []string{"us-east-1"}
	awsSessions := setupAwsSessions(configuration.Accounts, regions)
	awsSession := awsSessions["tuxlabs"]

	svc := lightsail.New(awsSession, &aws.Config{Region: aws.String(regions[0])})

	awsservices.CollectLightSailInstances(svc)

	if lib.CollectedData == nil {
		t.Error("Lightsail collected data is empty.")
	}
}
