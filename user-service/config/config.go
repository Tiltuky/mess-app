package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type AppConfig struct {
	HttpServer HttpServer `yaml:"httpServer"`
	GRPCServer GRPCServer `yaml:"grpcServer"`
}

type HttpServer struct {
	HttpHost string `yaml:"httpHost"`
	HttpPort string `yaml:"httpPort"`
}

type GRPCServer struct {
	GRPChost string `yaml:"grpcHost"`
	GRPCport string `yaml:"grpcPort"`
}

func LoadConfig(file string, cfg *AppConfig) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
