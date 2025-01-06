package push

import (
	"errors"
	"github.com/sirupsen/logrus"
)

type Push struct {
	Logger   *logrus.Logger
	apiToken string
}

func New(l *logrus.Logger, apiToken string) (*Push, error) {
	if apiToken == "" {
		return nil, errors.New("empty api token provided")
	}

	push := Push{
		Logger:   l,
		apiToken: apiToken,
	}

	return &push, nil
}
