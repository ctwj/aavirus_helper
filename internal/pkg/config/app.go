package config

import (
	"context"

	"github.com/spf13/viper"
)

type RunMode string

const (
	RunDebugMode   RunMode = "debug"
	RunReleaseMode RunMode = "release"
	RunTestMode    RunMode = "test"
)

var (
	App = struct {
		Mode RunMode `mapstructure:"mode"`
		Name string  `mapstructure:"name"`
	}{}
)

func LoadAppConfig(ctx context.Context, dir string) error {
	App.Mode = RunMode(viper.GetString("run_mode"))
	App.Name = viper.GetString("name")
	return nil
}
