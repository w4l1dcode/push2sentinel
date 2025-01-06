package sentinel

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Credentials struct {
	TenantID       string
	ClientID       string
	ClientSecret   string
	SubscriptionID string
}

type Sentinel struct {
	creds  Credentials
	logger *logrus.Logger

	azCreds    *azidentity.ClientSecretCredential
	httpClient *http.Client
}

func New(logger *logrus.Logger, creds Credentials) (*Sentinel, error) {
	sentinel := Sentinel{
		creds:  creds,
		logger: logger,
	}

	azCreds, err := azidentity.NewClientSecretCredential(creds.TenantID, creds.ClientID, creds.ClientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate to MS Sentinel: %v", err)
	}

	sentinel.azCreds = azCreds

	return &sentinel, nil
}
