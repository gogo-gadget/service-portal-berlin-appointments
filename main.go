package main

import (
	"github.com/gogo-gadget/service-portal-berlin-appointments/pkg/config"
	"github.com/gogo-gadget/service-portal-berlin-appointments/pkg/ticker"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.StandardLogger()

	configLoader := config.NewConfigLoader(config.Config{Logger: logger})

	cfg := &ticker.Config{}
	err := configLoader.LoadConfig(cfg)
	if err != nil {
		logger.WithError(err).Fatalf("could not load config")
	}

	t, err := ticker.NewTicker(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatalf("could not create ticker")
	}

	t.Run()
}
