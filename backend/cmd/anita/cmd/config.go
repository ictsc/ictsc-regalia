package cmd

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/server"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	// Dev 開発モード (default: false)
	Dev bool `json:"dev" mapstructure:"dev" yaml:"dev"`
	// Port サーバーがリッスンするポート番号 (default: 8080)
	Port int `json:"port" mapstructure:"port" yaml:"port"`

	// RDB RDB接続設定
	RDB struct {
		// Hostname ホスト名 (default: "mysql")
		Hostname string `json:"hostname" mapstructure:"hostname" yaml:"hostname"`
		// Port ポート番号 (default: 3306)
		Port int `json:"port" mapstructure:"port" yaml:"port"`
		// Username ユーザー名 (default: "root")
		Username string `json:"username" mapstructure:"username" yaml:"username"`
		// Password パスワード (default: "password")
		Password string `json:"password" mapstructure:"password" yaml:"password"`
		// Database データベース名 (default: "anita")
		Database string `json:"database" mapstructure:"database" yaml:"database"`
	} `json:"rdb" mapstructure:"rdb" yaml:"rdb"`
}

func provideRDBConfig(conf *Config) *bun.Config {
	return &bun.Config{
		Dev:      conf.Dev,
		Hostname: conf.RDB.Hostname,
		Port:     conf.RDB.Port,
		Username: conf.RDB.Username,
		Password: conf.RDB.Password,
		Database: conf.RDB.Database,
	}
}

func provideServerConfig(conf *Config) *server.Config {
	return &server.Config{
		Dev:  conf.Dev,
		Port: conf.Port,
	}
}

var (
	config     Config
	configFile string
)

func init() {
	viper.SetDefault("dev", false)
	viper.SetDefault("port", 8080)
	viper.SetDefault("rdb.hostname", "mysql")
	viper.SetDefault("rdb.port", 3306)
	viper.SetDefault("rdb.username", "root")
	viper.SetDefault("rdb.password", "password")
	viper.SetDefault("rdb.database", "anita")

	cobra.OnInitialize(func() {
		if len(configFile) > 0 {
			viper.SetConfigFile(configFile)
		} else {
			viper.AddConfigPath(".")
			viper.SetConfigName("config")
		}

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if errors.Is(err, viper.ConfigFileNotFoundError{}) {
				log.Panic(err)
			}
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Panic(err)
		}
	})
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print out the current config",
	Run: func(cmd *cobra.Command, args []string) {
		indented, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Panic(err)
		}

		log.Println("Printing out config\n" + string(indented))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path")
}
