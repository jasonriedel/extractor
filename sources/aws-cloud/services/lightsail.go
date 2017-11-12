package awsservices

import (
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/jasonriedel/extractor/lib"
)

func CollectLightSailInstances(svc *lightsail.Lightsail) {

	resp, err := svc.GetInstances(nil)
	if err != nil {
		lib.Log.Warning(err)
	}

	c := &lib.Collection{}
	for _,instance := range resp.Instances {
		uid := lib.UuidHash(*instance.Arn)
		c.StoreMap(uid, *instance)
	}
}
