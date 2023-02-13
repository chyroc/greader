package app

import (
	"fmt"

	"github.com/chyroc/go-loader"
)

type Config struct {
	MySQLHost     string `loader:"env,key=MYSQL_HOST"`
	MySQLUsername string `loader:"env,key=MYSQL_USERNAME"`
	MySQLPassword string `loader:"env,key=MYSQL_PASSWORD"`
	MySQLDatabase string `loader:"env,key=MYSQL_DATABASE"`
}

func loadConf() (*Config, error) {
	conf := new(Config)
	if err := loader.Load(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (conf *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.MySQLUsername, conf.MySQLPassword, conf.MySQLHost, conf.MySQLDatabase)
}
