package config

import (
	"context"

	"github.com/spf13/viper"
)

type EnvCfg struct {
	AndroidSdkPath string `mapstructure:"android_sdk_path"`
	AndroidSdkEnv  string `mapstructure:"android_sdk_env"`
}

var (
	Evn = EnvCfg{}
)

func LoadEnvConfig(ctx context.Context, workDir string) error {
	envCfg := viper.Sub("env")
	err := envCfg.Unmarshal(&Evn)
	if nil != err {
		return err
	}
	return nil
}
