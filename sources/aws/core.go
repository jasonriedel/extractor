package aws

import "github.com/aws/aws-sdk-go/aws/session"

func createAwsSession(account string) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Profile: account,
	})
}