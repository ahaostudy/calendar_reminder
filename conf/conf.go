package conf

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Server   Server   `yaml:"server"`
	MySQL    MySQL    `yaml:"mysql"`
	Redis    Redis    `yaml:"redis"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
	Email    Email    `yaml:"email"`
}

type Server struct {
	Address string `yaml:"address"`
	LogPath string `yaml:"log_path"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	DB       int    `yaml:"db"`
}

type RabbitMQ struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	VHost    string `yaml:"vhost"`
}

type Email struct {
	Addr   string `yaml:"addr"`
	Host   string `yaml:"host"`
	From   string `yaml:"from"`
	Email  string `yaml:"email"`
	Expire int    `yaml:"expire"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	var confPath string
	if gin.Mode() == gin.ReleaseMode {
		confPath = filepath.Join("conf", "conf_release.yaml")
	} else {
		confPath = filepath.Join("conf", "conf.yaml")
	}
	content, err := os.ReadFile(confPath)
	if err != nil {
		logrus.Fatalf("read conf file error - %v", err)
	}

	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		logrus.Fatalf("parse yaml error - %v", err)
	}
	if err := validator.Validate(conf); err != nil {
		logrus.Fatalf("validate conf error - %v", err)
	}

	bytes, _ := json.Marshal(conf)
	println(string(bytes))
}
