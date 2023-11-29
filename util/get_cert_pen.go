package util

import (
	"authentication/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
)

func GetPemCert(token *jwt.Token, auth0Domain string) (string, error) {
	cert := ""
	url := "https://" + auth0Domain + "/.well-known/jwks.json"
	resp, err := http.Get(url)
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = model.Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}
	if cert == "" {
		return cert, errors.New("unable to find appropriate key")
	}
	return cert, nil
}
