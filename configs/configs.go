package configs

import (
	"log"
	"strings"
	"sync"

	"github.com/chawin-a/wallet-monitor/internal/explorer"
	"github.com/chawin-a/wallet-monitor/internal/postgres"
	"github.com/chawin-a/wallet-monitor/internal/worker"
	"github.com/spf13/viper"
)

var (
	configOnce sync.Once
	config     *Config
)

type Config struct {
	Explorer explorer.Config `mapstructure:"explorer"`
	Postgres postgres.Config `mapstructure:"postgres"`
	Node     Node            `mapstructure:"node"`
	Worker   worker.Config   `mapstructure:"worker"`
	Wallets  []string        `mapstructure:"wallets"`
}

type Node struct {
	Url string `mapstructure:"url"`
}

func InitConfig() *Config {
	configOnce.Do(func() {
		configPath := "./configs"
		configName := "configs"
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if err := viper.ReadInConfig(); err != nil {
			log.Println("config file not found. using default/env config: " + err.Error())
		}
		viper.AutomaticEnv()

		viper.WatchConfig()
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	})
	return config
}
