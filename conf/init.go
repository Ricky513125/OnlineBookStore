package conf

import "github.com/spf13/viper"

var readErr error

func init() {
	viper.SetConfigFile("config.toml")
	viper.AddConfigPath("/etc/book-mall")
	viper.AddConfigPath("$HOME/.config/book-mall")
	viper.AddConfigPath(".")
	readErr = viper.ReadInConfig()
}
