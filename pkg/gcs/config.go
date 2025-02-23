package gcs

import (
	"encoding/base64"
	"errors"
	"net/url"
)

type Config struct {
	ServiceAccountFile string
	ServiceAccountJSON string
}

func (c Config) GetIngestrURI() (string, error) {
	missingCredentials := c.ServiceAccountFile == "" && c.ServiceAccountJSON == ""
	if missingCredentials {
		return "", errors.New("GCS: either service_account_file or service_account_json must be provided")
	}

	params := url.Values{}
	switch {
	case c.ServiceAccountFile != "":
		params.Set("credentials_path", c.ServiceAccountFile)
	case c.ServiceAccountJSON != "":
		params.Set(
			"credentials_base64",
			base64.StdEncoding.EncodeToString([]byte(c.ServiceAccountJSON)),
		)
	}
	uri := url.URL{
		Scheme:   "gs",
		RawQuery: params.Encode(),
	}
	return uri.String(), nil
}
