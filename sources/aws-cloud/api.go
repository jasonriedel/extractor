package awscloud

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jasonriedel/extractor/lib"
	)

func CreateAwsSession(account string) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Profile: account,
	})
}

func SetupAwsSessions(accountsMap map[string]string, regions []string) (map[string]*session.Session) {
	awsSessions := make(map[string]*session.Session, len(accountsMap))

	for accountName := range accountsMap {
		awsSession, err := CreateAwsSession(accountName)
		if err != nil {
			lib.Log.Critical(err)
		}

		awsSessions[accountName] = awsSession
	}
	return awsSessions
}
