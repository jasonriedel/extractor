package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"fmt"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/aws/aws-sdk-go/aws"
	"sync"
	"flag"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/op/go-logging"
	"github.com/jasonriedel/extractor/sources/aws-cloud/services"
	"github.com/jasonriedel/extractor/lib"
	"github.com/jasonriedel/extractor/sources/aws-cloud"
)

var wg sync.WaitGroup

var (
	fDaemonized = flag.Bool("d", false, "Run in daemon mode")
	fConfig = flag.String("config", "config/extractor.json", "Location of config file")
)

type ExtractionDetails struct {
	accountName string
	region string
}

func createChannelData(awsSessions map[string]*session.Session, accountsMap map[string]string, regions []string) (chan ExtractionDetails, map[string]*session.Session) {
	ch := make(chan ExtractionDetails)

	go func() {
		for accountName, _ := range accountsMap {
			for _,region := range regions {
				details := new(ExtractionDetails)
				details.accountName = accountName
				details.region = region
				ch <- *details
			}
		}
	}()

	return ch, awsSessions
}

func Collect(ch chan(ExtractionDetails), awsSessions map[string]*session.Session){
	details := <-ch
	accountName := details.accountName
	region := details.region
	awsSession := awsSessions[accountName]
	svc := lightsail.New(awsSession, &aws.Config{Region: aws.String(region)})
	msg := fmt.Sprintf("Collection Thread for %s, %s started", accountName, region)
	lib.Log.Info(msg)

	awsservices.CollectLightSailInstances(svc)
	fmt.Println(lib.CollectedData)
}


func getRegions() []string{
	var regions []string

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		//fmt.Println("Regions for", p.ID())
		if p.ID() != "aws-cloud-cn" && p.ID() != "us-gov-1" {
			for id, _ := range p.Regions() {
				regions = append(regions, id)
			}
		}
	}
	return regions
}


func main() {
	flag.Parse()

	//setup logging
	extractorlog1 := logging.NewLogBackend(os.Stderr, "", 0)
	extractorlog1Formatter := logging.NewBackendFormatter(extractorlog1, lib.LogFormat)
	logging.SetBackend(extractorlog1Formatter)

	configuration, err := lib.LoadConfiguration(*fConfig)
	if err != nil {
		lib.Log.Critical(err)
		os.Exit(1)
	}

	regions := getRegions()
	awsSessions := awscloud.SetupAwsSessions(configuration.Accounts, regions)

	ch, awsSessions := createChannelData(awsSessions, configuration.Accounts, regions)

	for i := 0; i < len(regions); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Collect(ch, awsSessions)
		}()
	}
	wg.Wait()
}