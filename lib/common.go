package lib

import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"encoding/json"
	"github.com/op/go-logging"
)

//logging setup
var Log = logging.MustGetLogger("logger")
var LogFormat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//global
var CollectedData = make(map[string]interface{})

type Collection []interface{}

type Configuration struct {
	Accounts map[string]string
}

func (c *Collection) StoreMap(uid string, data interface{}) {
	CollectedData[uid] = data
}

func UuidHash(arn string) string {
	return uuid.NewV5(uuid.NamespaceOID, arn).String()
}

func LoadConfiguration(path string) (config Configuration, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &config)
	return
}


