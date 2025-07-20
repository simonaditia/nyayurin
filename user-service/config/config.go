package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    int64  `json:"jwt_issuer"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Dbname    string `json:"dbname"`
	SSLMode   string `json:"sslmode"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type Config struct {
	App  App    `json:"app"`
	Psql PsqlDB `json:"db"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:      viper.GetString("APP_PORT"),
			AppEnv:       viper.GetString("APP_ENV"),
			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetInt64("JWT_ISSUER"),
		},
		Psql: PsqlDB{
			Host:      viper.GetString("DB_HOST"),
			Port:      viper.GetInt("DB_PORT"),
			User:      viper.GetString("DB_USER"),
			Password:  viper.GetString("DB_PASSWORD"),
			Dbname:    viper.GetString("DB_NAME"),
			SSLMode:   viper.GetString("DB_SSL_MODE"),
			DBMaxOpen: viper.GetInt("DB_MAX_OPEN"),
			DBMaxIdle: viper.GetInt("DB_MAX_IDLE"),
		},
	}
}
