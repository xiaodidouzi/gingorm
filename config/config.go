package config

import (
	"github.com/spf13/viper"
	"log"
)

// 配置结构体
type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port string `mapstructure:"port"`
	} `mapstructure:"app"`

	Database struct {
		Dsn          string `mapstructure:"dsn"`
		MaxIdleConns int    `mapstructure:"maxidleconns"`
		MaxOpenConns int    `mapstructure:"maxopenconns"`
	} `mapstructure:"database"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		DB       int    `mapstructure:"db"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`

	JWT struct {
		Secret string `mapstructure:"secret"`
		Expire int    `mapstructure:"expire"`
	} `mapstructure:"jwt"`
}

// 初始化配置
var AppConfig *Config

func InitConfig() {
	//
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file:%v", err)
	}
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct:%v", err)
	}
	InitDB()
	InitRedis()

}
