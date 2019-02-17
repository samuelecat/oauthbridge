// config
package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Version       string `yaml:"version"`
	HeadersPrefix string `yaml:"headers_prefix"`
	Providers     map[string]struct {
		BaseURI      string   `yaml:"base_uri"`
		ClientId     string   `yaml:"client_id"`
		ClientSecret string   `yaml:"client_secret"`
		Scopes       []string `yaml:"scopes"`
		TokenURL     string   `yaml:"token_url"`
		ExpireTime   int32    `yaml:"expire_time"`
	}
}

var Config configuration

func (c *configuration) getConf(file_path string) *configuration {
	yamlFile, err := ioutil.ReadFile(file_path)
	if err != nil {
		log.Fatalln("error loading configuration file: ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalln("error parsing configuration file: ", err)
	}
	return c
}

func loadConfig() {
	var file_path string
	if _, err := os.Stat("./conf/configuration.yml"); err == nil {
		// local path
		file_path = "./conf/configuration.yml"
	} else if _, err := os.Stat("/etc/oauthbridge/configuration.yml"); err == nil {
		file_path = "/etc/oauthbridge/configuration.yml"
	}

	if file_path == "" {
		log.Fatalln("error configuration file configuration.yml not found")
	}

	Config.getConf(file_path)
	if len(Config.Providers) == 0 {
		log.Fatalln("error parsing configuration file")
	}

	// check mandatory fields
	for name, provider := range Config.Providers {
		if provider.BaseURI == "" ||
			provider.ClientId == "" ||
			provider.ClientSecret == "" ||
			len(provider.Scopes) == 0 {
			log.Fatalln("error parsing configuration file: missing a mandatory field for provider", name)
		}
	}
}
