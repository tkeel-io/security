/*
Copyright 2021 The tKeel Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// DefaultConfigurationName is the default name of configuration.
	defaultConfigurationName = "auth"

	// DefaultConfigurationPath the default location of the configuration file.
	defaultConfigurationPath = "./"
)

type Config struct {
	Server   *ServerConfig   `json:"server" yaml:"server"`
	Database *DatabaseConfig `json:"database" yaml:"database"`
	Oauth2   *OAuth2Config   `json:"oauth2" yaml:"oauth2"`
	RBAC     *RBACConfig     `json:"rbac" yaml:"rbac"`
	Entity   *EntityConfig   `json:"entity" yaml:"entity" mapstruct:""`
}

type ServerConfig struct {
	Address string `json:"address" yaml:"address"`
}

type DatabaseConfig struct {
	Mysql *MysqlConf `json:"mysql" yaml:"mysql"`
}

type OAuth2Config struct {
	Redis          *RedisConf  `json:"redis" yaml:"redis"`
	AccessGenerate *AccessConf `json:"access_generate" yaml:"access_generate"`
}

type RBACConfig struct {
	Adapter *MysqlConf `json:"adapter" yaml:"adapter"`
}
type EntityConfig struct {
	SecurityKey string `json:"securitykey" yaml:"securitykey"`
}
type MysqlConf struct {
	DBName   string `json:"dbname" yaml:"dbname"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
}
type RedisConf struct {
	Addr string `json:"addr" yaml:"addr"`
	DB   int    `json:"db" yaml:"db"`
}
type AccessConf struct {
	SecurityKey string `json:"security_key" yaml:"security_key"`
}

func New() *Config {
	return &Config{
		Server: &ServerConfig{
			Address: ":8080",
		},
	}
}

// TryLoadFromDisk loads configuration from default location after server startup.
// return nil error if configuration file not exists.
func TryLoadFromDisk() (*Config, error) {
	viper.SetConfigName(defaultConfigurationName)
	viper.AddConfigPath(defaultConfigurationPath)

	// Load from current working directory, only used for debugging.
	viper.AddConfigPath(".")

	// Load from Environment variables.
	viper.SetEnvPrefix("auth")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("viper read in config %w", err)
		}
		return nil, fmt.Errorf("error parsing configuration file %w", err)
	}

	conf := New()

	if err := viper.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("unmarshal conf %w", err)
	}

	return conf, nil
}
