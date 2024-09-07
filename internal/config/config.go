package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DBConfig `yaml:"db"`
	Server   Server   `yaml:"server"`
	//Broker   Broker   `yaml:broker` //TODO KAFKA??? or RMQ???
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	InMemory bool   `yaml:"inMemory"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Broker struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	UserName     string `yaml:"user"`
	Password     string `yaml:"password"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchangeType"`
	RoutingKey   string `yaml:"routingKey"`
}

func GetBannersConfig(configFilePath string) Config {
	conf := &Config{}

	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(conf); err != nil {
		panic(err)
	}

	return *conf
}
