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
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const (
	// DefaultConfigurationName is the default name of configuration
	defaultConfigurationName = "auth"

	// DefaultConfigurationPath the default location of the configuration file
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

// Package config saves configuration for running KubeSphere components
//
// Config can be configured from command line flags and configuration file.
// Command line flags hold higher priority than configuration file. But if
// component Endpoint/Host/APIServer was left empty, all of that component
// command line flags will be ignored, use configuration file instead.
// For example, we have configuration file
//
// mysql:
//   host: mysql.kubesphere-system.svc
//   username: root
//   password: password
//
// At the same time, have command line flags like following:
//
// --mysql-host mysql.openpitrix-system.svc --mysql-username king --mysql-password 1234
//
// We will use `king:1234@mysql.openpitrix-system.svc` from command line flags rather
// than `root:password@mysql.kubesphere-system.svc` from configuration file,
// cause command line has higher priority. But if command line flags like following:
//
// --mysql-username root --mysql-password password
//
// we will `root:password@mysql.kubesphere-system.svc` as input, cause
// mysql-host is missing in command line flags, all other mysql command line flags
// will be ignored.

//// Config defines everything needed for apiserver to deal with external services
//type Config struct {
//	DevopsOptions         *jenkins.Options                           `json:"devops,omitempty" yaml:"devops,omitempty" mapstructure:"devops"`
//	SonarQubeOptions      *sonarqube.Options                         `json:"sonarqube,omitempty" yaml:"sonarQube,omitempty" mapstructure:"sonarqube"`
//	KubernetesOptions     *k8s.KubernetesOptions                     `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty" mapstructure:"kubernetes"`
//	ServiceMeshOptions    *servicemesh.Options                       `json:"servicemesh,omitempty" yaml:"servicemesh,omitempty" mapstructure:"servicemesh"`
//	NetworkOptions        *network.Options                           `json:"network,omitempty" yaml:"network,omitempty" mapstructure:"network"`
//	LdapOptions           *ldap.Options                              `json:"-,omitempty" yaml:"ldap,omitempty" mapstructure:"ldap"`
//	RedisOptions          *cache.Options                             `json:"redis,omitempty" yaml:"redis,omitempty" mapstructure:"redis"`
//	S3Options             *s3.Options                                `json:"s3,omitempty" yaml:"s3,omitempty" mapstructure:"s3"`
//	OpenPitrixOptions     *openpitrix.Options                        `json:"openpitrix,omitempty" yaml:"openpitrix,omitempty" mapstructure:"openpitrix"`
//	MonitoringOptions     *prometheus.Options                        `json:"monitoring,omitempty" yaml:"monitoring,omitempty" mapstructure:"monitoring"`
//	LoggingOptions        *logging.Options                           `json:"logging,omitempty" yaml:"logging,omitempty" mapstructure:"logging"`
//	AuthenticationOptions *authoptions.AuthenticationOptions         `json:"authentication,omitempty" yaml:"authentication,omitempty" mapstructure:"authentication"`
//	AuthorizationOptions  *authorizationoptions.AuthorizationOptions `json:"authorization,omitempty" yaml:"authorization,omitempty" mapstructure:"authorization"`
//	MultiClusterOptions   *multicluster.Options                      `json:"multicluster,omitempty" yaml:"multicluster,omitempty" mapstructure:"multicluster"`
//	EventsOptions         *events.Options                            `json:"events,omitempty" yaml:"events,omitempty" mapstructure:"events"`
//	AuditingOptions       *auditing.Options                          `json:"auditing,omitempty" yaml:"auditing,omitempty" mapstructure:"auditing"`
//	AlertingOptions       *alerting.Options                          `json:"alerting,omitempty" yaml:"alerting,omitempty" mapstructure:"alerting"`
//	NotificationOptions   *notification.Options                      `json:"notification,omitempty" yaml:"notification,omitempty" mapstructure:"notification"`
//	KubeEdgeOptions       *kubeedge.Options                          `json:"kubeedge,omitempty" yaml:"kubeedge,omitempty" mapstructure:"kubeedge"`
//	MeteringOptions       *metering.Options                          `json:"metering,omitempty" yaml:"metering,omitempty" mapstructure:"metering"`
//	GatewayOptions        *gateway.Options                           `json:"gateway,omitempty" yaml:"gateway,omitempty" mapstructure:"gateway"`
//}

// TryLoadFromDisk loads configuration from default location after server startup
// return nil error if configuration file not exists
func TryLoadFromDisk() (*Config, error) {
	viper.SetConfigName(defaultConfigurationName)
	viper.AddConfigPath(defaultConfigurationPath)

	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")

	// Load from Environment variables
	viper.SetEnvPrefix("auth")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		} else {
			return nil, fmt.Errorf("error parsing configuration file %s", err)
		}
	}

	conf := New()

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
