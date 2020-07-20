package model

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServiceConfiguration struct {
	Name        string `yaml:"name"`
	Image       string `yaml:"image"`
	ExposedPort string `yaml:"exposedPort"`
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

type auth struct {
	BaUser     string `json:"baUser"`
	BaPassword string `json:"baPassword"`
}

type ServiceParams struct {
	ServiceInstanceName string `json:"serviceInstanceName"`
	Auth                auth   `json:"auth"`
}

func Marshal(rawJson []byte) (*ServiceParams, error) {
	var params ServiceParams
	if err := json.Unmarshal(rawJson, &params); err != nil {
		return nil, err
	}
	return &params, nil
}
