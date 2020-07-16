package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServiceConfiguration struct {
	Name        string `yaml:"name"`
	Image       string `yaml:"image"`
	Port        string `yaml:"port"`
	Description string `yaml:"description"`
	ServiceId   string `yaml:"serviceId"`
	PlanId      string `yaml:"planId"`
}

type Services struct {
	Catalog []ServiceConfiguration `yaml:"catalog"`
}

func Parse(filePath string) (Services, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Services{}, err
	}

	var services Services
	if err := yaml.Unmarshal(fileContents, &services); err != nil {
		return Services{}, err
	}
	return services, nil
}
