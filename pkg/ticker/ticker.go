package ticker

import (
	"time"

	"github.com/gogo-gadget/service-portal-berlin-appointments/pkg/checker"
	"github.com/gogo-gadget/service-portal-berlin-appointments/pkg/notifier"
	"github.com/sirupsen/logrus"
)

type Config struct {
	CheckIntervalInSeconds  int
	RepeatIntervalInSeconds int

	CheckerCfg  checker.Config
	NotifierCfg notifier.SMTPNotifierConfig
}

type Ticker interface {
	Run()
}

type ticker struct {
	cfg *Config
	t   *time.Ticker

	checker  checker.Checker
	notifier notifier.Notifier

	logger *logrus.Logger
}

func NewTicker(cfg *Config, logger *logrus.Logger) (Ticker, error) {
	c, err := checker.NewChecker(cfg.CheckerCfg, logger)
	if err != nil {
		return nil, err
	}
	n := notifier.NewSMTPNotifier(cfg.NotifierCfg, logger)

	t := &ticker{
		cfg:      cfg,
		checker:  c,
		notifier: n,
		logger:   logger,
	}

	return t, nil
}

func (t *ticker) Run() {
	logger := t.logger

	timeTicker := time.NewTicker(time.Duration(t.cfg.CheckIntervalInSeconds) * time.Second)
	t.t = timeTicker

	for {
		<-t.t.C
		appointments, err := t.checker.GetAppointments()
		if err != nil {
			// skip to next iteration
			logger.WithError(err).Errorf("could not check appointments")
			continue
		}

		if len(appointments) == 0 {
			// nothing to do be done skip to next iteration
			logger.Infof("no appointments received")
			continue
		}

		err = t.notifier.Notify(appointments)
		if err != nil {
			logger.WithError(err).Errorf("could not notify about appointments")
		}

		// Wait before sending another potential notification
		<-time.After(time.Duration(t.cfg.RepeatIntervalInSeconds) * time.Second)
	}
}
