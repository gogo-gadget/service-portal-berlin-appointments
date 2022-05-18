package checker

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type Config struct {
	CookieURL      string
	AppointmentURL string
}

type Appointment struct {
	Day   string
	Month string
}

type Checker interface {
	GetAppointments() ([]Appointment, error)
}

type checker struct {
	cfg Config

	cookieURL      *url.URL
	appointmentURL *url.URL

	client *http.Client

	logger *logrus.Logger
}

func NewChecker(cfg Config, logger *logrus.Logger) (Checker, error) {
	// create http client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	cookieURL, err := url.Parse(cfg.CookieURL)
	if err != nil {
		logger.WithError(err).Error("could not parse cookie url")
		return nil, err
	}

	appointmentURL, err := url.Parse(cfg.AppointmentURL)
	if err != nil {
		logger.WithError(err).Error("could not parse appointment url")
		return nil, err
	}

	c := &checker{
		cfg:            cfg,
		cookieURL:      cookieURL,
		appointmentURL: appointmentURL,
		client:         client,
		logger:         logger,
	}

	return c, nil
}

func (c *checker) GetAppointments() ([]Appointment, error) {
	logger := c.logger

	// cookie request
	cookieRequest, err := http.NewRequest(http.MethodGet, c.cookieURL.String(), nil)
	if err != nil {
		logger.WithError(err).Error("could not create cookie request")
		return nil, err
	}
	cookieRequest.Header.Set("host", c.cookieURL.Host)

	cookieResp, err := c.client.Do(cookieRequest)
	if err != nil {
		logger.WithError(err).Error("could not send http cookieRequest")
		return nil, err
	}

	cookie := cookieResp.Header.Get("Set-Cookie")
	if cookie == "" {
		errMsg := "could not get cookie from http response"
		logger.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// appointment request
	appointmentRequest, err := http.NewRequest(http.MethodGet, c.appointmentURL.String(), nil)
	if err != nil {
		logger.WithError(err).Error("could not create appointment request")
		return nil, err
	}
	appointmentRequest.Header.Set("host", c.appointmentURL.Host)
	appointmentRequest.Header.Set("cookie", cookie)

	appointmentResp, err := c.client.Do(appointmentRequest)

	if appointmentResp.StatusCode < 200 || appointmentResp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("received non okay status code: %v", cookieResp.StatusCode)
		logger.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	doc, err := goquery.NewDocumentFromReader(appointmentResp.Body)
	if err != nil {
		logger.WithError(err).Error("could not doc from appointment response body")
		return nil, err
	}

	appointments := []Appointment{}

	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		month := ""
		table.Find("th").Each(func(i int, th *goquery.Selection) {
			if th.HasClass("month") {
				month = th.Text()
			}
		})
		table.Find("td").Each(func(i int, td *goquery.Selection) {
			if td.HasClass("buchbar") {
				day := td.Text()
				appointment := Appointment{
					Day:   day,
					Month: month,
				}
				appointments = append(appointments, appointment)
			}
		})
	})

	return appointments, nil
}
