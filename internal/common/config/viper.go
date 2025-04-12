package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	viperOnce sync.Once
	viperErr  error
)

func init() {
	viperOnce.Do(func() {
		viperErr = newViperConfig()
	})
	if viperErr != nil {
		panic(viperErr.Error())
	}
}

// NewViperConfig initializes the Viper configuration.
func newViperConfig() error {
	relPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}

	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	_ = viper.BindEnv("endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	return viper.ReadInConfig()
}

func getRelativePathFromCaller() (string, error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	_, here, _, _ := runtime.Caller(0)
	relPath, err := filepath.Rel(callerPwd, filepath.Dir(here))
	//fmt.Printf("callerPwd: %s, viper config relative path: %s\n", callerPwd, relPath)
	return relPath, err
}
