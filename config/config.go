package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const dopplerBaseURI = "https://api.doppler.com/v3/configs/config/secrets"

var (
	dopplerToken = os.Getenv("DOP_TOKEN")
	environment  = os.Getenv("RUN_MODE")
	envName      = os.Getenv("DOPPLER_ENV_NAME")
)

type Config struct {
	Project   string
	Secrets   map[string]string
	serverURI string
	l         sync.RWMutex
}

func Configure(prefix string) *Config {
	if dopplerToken == "" {
		log.Println("doppler token is not set")
	}
	c := Config{
		Project:   prefix,
		serverURI: dopplerBaseURI,
	}
	err := c.GetDopplerSecrets()
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

func (c *Config) GetDopplerSecrets() error {
	secrets := make(map[string]string)
	url := fmt.Sprintf("%s?project=better_secrets&config=dev&include_dynamic_secrets=false&include_managed_secrets=true",
		c.serverURI,
	)

	client := http.DefaultClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("accepts", "application/json")
	req.Header.Add("Authorization", "Bearer "+dopplerToken)

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
