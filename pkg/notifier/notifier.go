package notifier

import (
	"fmt"
	"net/smtp"

	"github.com/gogo-gadget/service-portal-berlin-appointments/pkg/checker"
	"github.com/sirupsen/logrus"
)

type SMTPNotifierConfig struct {
	Address  string
	Identity string
	Host     string

	From    string
	To      []string
	Subject string

	Username string
	Password string

	MsgURL string
}

type Notifier interface {
	Notify(appointments []checker.Appointment) error
}

type smtpNotifier struct {
	cfg    SMTPNotifierConfig
	logger *logrus.Logger
}

func NewSMTPNotifier(cfg SMTPNotifierConfig, logger *logrus.Logger) Notifier {
	n := &smtpNotifier{
		cfg:    cfg,
		logger: logger,
	}

	return n
}

func (n *smtpNotifier) Notify(appointments []checker.Appointment) error {
	logger := n.logger
	logger.Infof("notifying about %v appointments", len(appointments))

	msgHeader := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n", n.cfg.From, n.cfg.To, n.cfg.Subject)
	msgBody := fmt.Sprintf("Hello!\nThere are %d new apppointments on %s!", len(appointments), n.cfg.MsgURL)

	msg := msgHeader + msgBody

	err := smtp.SendMail(n.cfg.Address,
		smtp.PlainAuth(n.cfg.Identity, n.cfg.Username, n.cfg.Password, n.cfg.Host),
		n.cfg.From,
		n.cfg.To,
		[]byte(msg))

	return err
}
