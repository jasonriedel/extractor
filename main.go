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
	"github.com/tuxninja/extractor/aws"
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
	log.Info(msg)

	CollectLightSailInstances(svc)
	fmt.Println(CollectedData)
}



func setupAwsSessions(accountsMap map[string]string, regions []string) (map[string]*session.Session) {
	awsSessions := make(map[string]*session.Session, len(accountsMap))

	for accountName := range accountsMap {
		awsSession, err := createAwsSession(accountName)
		if err != nil {
			log.Critical(err)
		}

		awsSessions[accountName] = awsSession
	}
	return awsSessions
}


func getRegions() []string{
	var regions []string

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		//fmt.Println("Regions for", p.ID())
		if p.ID() != "aws-cn" && p.ID() != "us-gov-1" {
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
	extractorlog1Formatter := logging.NewBackendFormatter(extractorlog1, logFormat)
	logging.SetBackend(extractorlog1Formatter)

	configuration, err := loadConfiguration(*fConfig)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}

	regions := getRegions()
	awsSessions := setupAwsSessions(configuration.Accounts, regions)

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
