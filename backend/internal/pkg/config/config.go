package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// DBConfig はデータベース接続情報を保持します。
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN はデータベース接続文字列 (DSN) を生成します。
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// LoadDBConfig は指定されたパスとファイル名からDB設定を読み込みます。
func LoadDBConfig(path string, fileName string) (config DBConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.UnmarshalKey("database", &config)
	return
}
