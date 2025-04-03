package config

import (
	"github.com/spf13/viper"
	"strings"
)

// NewViperConfig initializes the Viper configuration.
func NewViperConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	_ = viper.BindEnv("endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")

	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
