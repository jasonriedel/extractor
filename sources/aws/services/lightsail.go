package aws

import (
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/op/go-logging"
)

//logging setup
var log = logging.MustGetLogger("logger")
var logFormat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func CollectLightSailInstances(svc *lightsail.Lightsail) {

	resp, err := svc.GetInstances(nil)
	if err != nil {
		log.Warning(err)
	}
	c := &Collection{}
	for _,instance := range resp.Instances {
		uid := uuidHash(*instance.Arn)
		c.StoreMap(uid, *instance)
	}
}
