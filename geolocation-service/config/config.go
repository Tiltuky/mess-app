package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Local     Local     `yaml:"local"`
	DB        DBConfig  `yaml:"db"`
	RedisConf RedisConf `yaml:"redis"`
	Kafka     Kafka     `yaml:"kafka"`
	Payments  Payments  `yaml:"payments"`
}

type Local struct {
	Port int `yaml:"port"`
}

type RedisConf struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	DBName   string        `yaml:"dbname"`
	SSlMode  string        `yaml:"sslmode"`
	Driver   string        `yaml:"driver"`
	TimeOut  time.Duration `yaml:"timeout"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type Payments struct {
	PublishableKey  string `yaml:"publishable_key"`
	SecretKey       string `yaml:"secret_key"`
	ProductID       string `yaml:"product_id"`
	PaymentMethodID string `yaml:"payment_method_id"`
	TrialEnd        int64  `yaml:"trial_end"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	configPath, _ := os.LookupEnv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Путь до конфига не найден в енв файле")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Файл конфига не найден")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Не можем прочитать конфиг %s", err)
	}
	log.Printf("Config: %+v", cfg)
	return &cfg
}
