package config

type Config struct {
	AppName     string      `mapstructure:"appname"`
	CacheConfig CacheConfig `mapstructure:"cacheconfig"`
	LogConfig   LogConfig   `mapstructure:"logconfig"`
}

type CacheConfig struct {
	CacheHost     string `mapstructure:"cachehost"`
	CachePort     int    `mapstructure:"cacheport"`
	CacheUsername string `mapstructure:"cacheusername"`
	CachePassword string `mapstructure:"cachepassword"`
	CacheDB       int    `mapstructure:"cachedb"`
	CacheTTL      int    `mapstructure:"cachettl"`
	CacheEnabled  bool   `mapstructure:"cacheenabled"`
}

type LogConfig struct {
	LogAppName string `mapstructure:"logappname"`
	LogEnv     string `mapstructure:"logenv"`
}
