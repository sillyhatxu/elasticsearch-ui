package main

import (
	"flag"
	"fmt"
	"github.com/sillyhatxu/elasticsearch-ui/api"
	"github.com/sillyhatxu/elasticsearch-ui/config"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	flag.StringVar(&config.Conf.ServerHost, "p", "0.0.0.0:8080", "local server address")
	flag.StringVar(&config.Conf.URL, "url", "http://127.0.0.1:9200", "es address")
	flag.Usage = func() {
		fmt.Println(fmt.Sprintf("Usage of %s:\n", os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	err := api.InitialAPI()
	if err != nil {
		logrus.Errorf("elasticsearch-ui error. Error : %v", err)
	} else {
		logrus.Info("project close")
	}
}
