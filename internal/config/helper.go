package config

import (
	"github.com/spf13/viper"
)

func configDefault() Config {
	return Config{
		AppName: "Api-Anime",
		CacheConfig: CacheConfig{
			CacheHost:    "localhost",
			CachePort:    6379,
			CacheDB:      0,
			CacheTTL:     60,
			CacheEnabled: false,
		},
		LogConfig: LogConfig{
			LogAppName: "Api-Anime",
			LogEnv:     "production",
		},
	}
}

func setConfigDefault() {
	defaults := configDefault()

	viper.SetDefault("appname", defaults.AppName)
	viper.SetDefault("cacheconfig.cachehost", defaults.CacheConfig.CacheHost)
	viper.SetDefault("cacheconfig.cacheport", defaults.CacheConfig.CachePort)
	viper.SetDefault("cacheconfig.cacheusername", defaults.CacheConfig.CacheUsername)
	viper.SetDefault("cacheconfig.cachepassword", defaults.CacheConfig.CachePassword)
	viper.SetDefault("cacheconfig.cachedb", defaults.CacheConfig.CacheDB)
	viper.SetDefault("cacheconfig.cachettl", defaults.CacheConfig.CacheTTL)
	viper.SetDefault("cacheconfig.enablecache", defaults.CacheConfig.CacheEnabled)
	viper.SetDefault("logconfig.logappname", defaults.LogConfig.LogAppName)
	viper.SetDefault("logconfig.logenv", defaults.LogConfig.LogEnv)
}
