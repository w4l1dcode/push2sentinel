package main

import (
	"flag"
	_ "fmt"
	"github.com/sirupsen/logrus"
	"github.com/w4l1dcode/push2sentinel/config"
	"github.com/w4l1dcode/push2sentinel/pkg/push"
	_ "io/ioutil"
	_ "net/http"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	// ---

	pushClient, err := push.New(logger, conf.Push.ApiToken)
	if err != nil {
		logger.WithError(err).Fatal("failed to create aikido client")
	}

	// --

	accounts, err := pushClient.GetAccounts(1)
	if err != nil {
		logger.WithError(err).Fatal("failed to get issues")
	}
	logger.Infof("Retrieved %d accounts:", len(accounts))
	logger.Println(accounts)

	// --

	apps, err := pushClient.GetApps(1)
	if err != nil {
		logger.WithError(err).Fatal("failed to get apps")
	}
	logger.Infof("Retrieved %d apps:", len(apps))
	logger.Println(apps)

	// --

	browsers, err := pushClient.GetBrowsers(1)
	if err != nil {
		logger.WithError(err).Fatal("failed to get browsers")
	}
	logger.Infof("Retrieved %d browsers:", len(browsers))
	logger.Println(browsers)

	// --

	employees, err := pushClient.GetEmployees(1)
	if err != nil {
		logger.WithError(err).Fatal("failed to get employees")
	}
	logger.Infof("Retrieved %d employees:", len(employees))
	logger.Println(employees)

	// --

	findings, err := pushClient.GetFindings(1)
	if err != nil {
		logger.WithError(err).Fatal("failed to get findings")
	}
	logger.Infof("Retrieved %d findings:", len(findings))
	logger.Println(findings)

}
