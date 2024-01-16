package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/ctwj/aavirus_helper/internal/lib"
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

var (
	once sync.Once
	//confLoader *configLoader
)

type ConfigModuleLoader func(ctx context.Context, workDir string) error

func Load(ctx context.Context, cfgPath string, moduleLoaders ...ConfigModuleLoader) (err error) {
	once.Do(func() {
		c := &configLoader{
			cfgPath: cfgPath,
			loaders: append([]ConfigModuleLoader{
				LoadAppConfig,
				LoadLogConfig,
				LoadEnvConfig,
				// loadSettingConfig,
				// LoadRemoteConfig,
			}, moduleLoaders...),
		}
		if er := c.init(ctx); er != nil {
			err = fmt.Errorf("config init error:%w", er)
			return
		}

		c.watch(ctx)
		//confLoader = c
	})
	return err
}

type configLoader struct {
	loaders []ConfigModuleLoader
	cfgPath string
}

func (c *configLoader) init(ctx context.Context) error {
	if len(c.cfgPath) > 0 {
		viper.SetConfigFile(c.cfgPath)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return c.load(ctx)
}

func (c *configLoader) load(ctx context.Context) error {

	dir := os.Getenv(EnvAppDirName)
	if len(dir) == 0 {
		_dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("os.Getwd error:%w", err)
		}
		dir = _dir
	}

	AppDir = dir
	initDir()

	for _, loader := range c.loaders {
		if err := loader(ctx, dir); err != nil {
			return err
		}
	}

	return nil
}

func initDir() {
	TmpDir = path.Join(AppDir, "tmp")
	WorkDir = path.Join(AppDir, "workspace")
	OutputDir = path.Join(AppDir, "output")
	err := lib.CreateDirectoryIfNotExists(TmpDir)
	if err != nil {
		fmt.Println(ERROR_TAG, "Error creating conf directory:", err)
	}
	lib.CreateDirectoryIfNotExists(WorkDir)
	if err != nil {
		fmt.Println(ERROR_TAG, "Error creating work directory:", err)
	}
	lib.CreateDirectoryIfNotExists(OutputDir)
	if err != nil {
		fmt.Println(ERROR_TAG, "Error creating output directory:", err)
	}
}

func (c *configLoader) watch(ctx context.Context) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
