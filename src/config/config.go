package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type AppConfig struct {
	Mysql    MysqlConfig
	Mode     string
	Gorm     GormConfig
	Port     int
	Redis    RedisConfig
	LogLevel string
	Es       EsConfig
}

type MysqlConfig struct {
	DB   string
	Host string
	User string
	Pwd  string
}

type GormConfig struct {
	Debug bool
}

type RedisConfig struct {
	Host     string
	Port     int
	PoolSize int
}

type EsConfig struct {
	Host     string
	Username string
	Password string
}

var Config AppConfig

func InitConfig() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yml", "配置文件路径")
	flag.Parse()
	ParseConfig(configPath)
}

func ParseConfig(configPath string) {
	bArr, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("can not find config file [%s]! \n %s", configPath, err))
	}

	if err = yaml.Unmarshal(bArr, &Config); err != nil {
		panic(fmt.Sprintf("parse config file [%s] fail! \n %s", configPath, err))
	}
}
