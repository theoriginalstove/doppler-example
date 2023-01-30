package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const dopplerBaseURI = "https://api.doppler.com/v3/configs/"

var (
	dopplerToken = os.Getenv("DOP_TOKEN")
	environment  = os.Getenv("RUNMODE")
)

type Config struct {
	Project   string
	Secrets   map[string]string
	serverURI string
}

func Configure(prefix string) *Config {
	c := Config{
		Project:   prefix,
		serverURI: dopplerBaseURI,
	}
	err := c.getDopplerSecrets()
	if err != nil {
		log.Fatal("unable to retrieve and set secrets from doppler")
	}
	return &c
}

type dopplerSecretKey struct {
	Raw      string `json:"raw"`
	Computed string `json:"computed"`
}

type dopplerSecretsRes struct {
	Secrets map[string]dopplerSecretKey `json:"secrets"`
}

func (c *Config) getDopplerSecrets() error {
	secrets := make(map[string]string)
	url := fmt.Sprintf("%s?project=%v&config=%v&include_dynamic_secrets=true&dynamic_secrets_ttl_sec=1800",
		c.serverURI,
		c.Project,
		environment,
	)

	client := http.DefaultClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(dopplerToken, "")
	req.Header.Add("accept", "application/json")
	req.Header.Add("accepts", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var doppSecretRes dopplerSecretsRes
	err = json.Unmarshal(body, &doppSecretRes)
	if err != nil {
		return err
	}
	for k, v := range doppSecretRes.Secrets {
		secrets[k] = v.Computed
	}
	c.Secrets = secrets
	return nil
}
