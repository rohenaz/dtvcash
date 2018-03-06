package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	EnvMysqlHost = "MYSQL_HOST"
	EnvMysqlUser = "MYSQL_USER"
	EnvMysqlPass = "MYSQL_PASS"
	EnvMysqlDb   = "MYSQL_DB"
)

type MysqlConfig struct {
	Host     string
	Username string
	Password string
	Database string
}

func init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.memo")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found")
	}
}

func GetMysqlConfig() MysqlConfig {
	return MysqlConfig{
		Host:     viper.GetString(EnvMysqlHost),
		Username: viper.GetString(EnvMysqlUser),
		Password: viper.GetString(EnvMysqlPass),
		Database: viper.GetString(EnvMysqlDb),
	}
}
