package requests

import (
	json2 "atlas-party/json"
	"atlas-party/retry"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	BaseRequest string = "http://atlas-nginx:80"
)

type configuration struct {
	retries int
}

type Configurator func(c *configuration)

func SetRetries(amount int) Configurator {
	return func(c *configuration) {
		c.retries = amount
	}
}

func Get(l logrus.FieldLogger) func(url string, resp interface{}, configurators ...Configurator) error {
	return func(url string, resp interface{}, configurators ...Configurator) error {
		c := &configuration{retries: 1}
		for _, configurator := range configurators {
			configurator(c)
		}

		var r *http.Response
		get := func(attempt int) (bool, error) {
			var err error
			r, err = http.Get(url)
			if err != nil {
				l.Warnf("Failed calling GET on %s, will retry.", url)
				return true, err
			}
			return false, nil
		}
		err := retry.Try(get, c.retries)
		if err != nil {
			l.WithError(err).Errorf("Unable to successfully call GET on %s.", url)
			return err
		}
		err = ProcessResponse(r, resp)
		return err
	}
}

func ProcessResponse(r *http.Response, rb interface{}) error {
	err := json2.FromJSON(rb, r.Body)
	if err != nil {
		return err
	}

	return nil
}
