package main

import (
	"context"
	"flag"
	_ "fmt"
	"github.com/sirupsen/logrus"
	"github.com/w4l1dcode/push2sentinel/config"
	"github.com/w4l1dcode/push2sentinel/pkg/push"
	msSentinel "github.com/w4l1dcode/push2sentinel/pkg/sentinel"
	_ "io/ioutil"
	_ "net/http"
	"sync"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "push2sentinel_config.yml", "The YAML configuration file.")
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
	logger.WithField("level", logrusLevel.String()).Info("set log level")
	logger.SetLevel(logrusLevel)

	// ---

	errors := make(chan error)
	wg := &sync.WaitGroup{}

	// Declare variables to hold the results for each type
	var accounts []map[string]string
	var apps []map[string]string
	var browsers []map[string]string
	var employees []map[string]string
	var findings []map[string]string

	// Fetch accounts
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Retrieving accounts")
		client, err := push.New(logger, conf.Push.ApiToken)
		if err != nil {
			errors <- err
			return
		}
		accounts, err = client.GetAccounts(conf.Push.LookbackHours)
		if err != nil {
			errors <- err
		}
	}()

	// Fetch apps
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Retrieving apps")
		client, err := push.New(logger, conf.Push.ApiToken)
		if err != nil {
			errors <- err
			return
		}
		apps, err = client.GetApps(conf.Push.LookbackHours)
		if err != nil {
			errors <- err
		}
	}()

	// Fetch browsers
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Retrieving browsers")
		client, err := push.New(logger, conf.Push.ApiToken)
		if err != nil {
			errors <- err
			return
		}
		browsers, err = client.GetBrowsers(conf.Push.LookbackHours)
		if err != nil {
			errors <- err
		}
	}()

	// Fetch employees
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Retrieving employees")
		client, err := push.New(logger, conf.Push.ApiToken)
		if err != nil {
			errors <- err
			return
		}
		employees, err = client.GetEmployees(conf.Push.LookbackHours)
		if err != nil {
			errors <- err
		}
	}()

	// Fetch findings
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Retrieving findings")
		client, err := push.New(logger, conf.Push.ApiToken)
		if err != nil {
			errors <- err
			return
		}
		findings, err = client.GetFindings(conf.Push.LookbackHours)
		if err != nil {
			errors <- err
		}
	}()

	// ---

	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	logger.Info("waiting for log ingestion to finish")
	select {
	case err := <-errors:
		logger.WithError(err).Fatal("failed to retrieve logs")
	case <-doneChan:
		logger.Info("finished retrieving logs")
	}

	// ---

	sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
		TenantID:       conf.Microsoft.TenantID,
		ClientID:       conf.Microsoft.AppID,
		ClientSecret:   conf.Microsoft.SecretKey,
		SubscriptionID: conf.Microsoft.SubscriptionID,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not create MS Sentinel client")
	}

	// ---

	totalLength := len(accounts) + len(apps) + len(employees) + len(browsers) + len(findings)
	allLogs := make([]map[string]string, 0, totalLength)
	allLogs = append(allLogs, accounts...)
	allLogs = append(allLogs, apps...)
	allLogs = append(allLogs, employees...)
	allLogs = append(allLogs, browsers...)
	allLogs = append(allLogs, findings...)

	// ---

	logger.WithField("total", len(allLogs)).Info("shipping off push security logs to Sentinel")

	if err := sentinel.SendLogs(ctx, logger,
		conf.Microsoft.DataCollection.Endpoint,
		conf.Microsoft.DataCollection.RuleID,
		conf.Microsoft.DataCollection.StreamName,
		allLogs); err != nil {
		logger.WithError(err).Fatal("could not ship logs to sentinel")
	}

	logger.WithField("total", len(allLogs)).Info("successfully sent logs to sentinel")
}
