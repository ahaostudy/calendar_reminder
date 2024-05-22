package conf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"

	"gopkg.in/validator.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Server Server `yaml:"server"`
	MySQL  MySQL  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
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

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	confPath := filepath.Join("conf", "conf.yaml")
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
