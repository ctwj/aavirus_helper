package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type EnvCfg struct {
	AndroidSdkPath          string `mapstructure:"android_sdk_path"`
	AndroidSdkEnv           string `mapstructure:"android_sdk_env"`
	AndroidBuildToolPath    string `mapstructure:"android_build_tool_path"`
	AndroidBuildToolVersion string `mapstructure:"android_build_tool_version"`
	ApkToolName             string `mapstructure:"apk_tool_name"`
	ToolDir                 string `mapstructure:"tool_dir"`
}

var (
	Env = EnvCfg{}
)

func LoadEnvConfig(ctx context.Context, workDir string) error {
	envCfg := viper.Sub("env")
	err := envCfg.Unmarshal(&Env)

	fmt.Println("AndroidSdkPath", Env.AndroidSdkPath)

	if Env.AndroidSdkPath == "" {
		// 加载环境变量
		Env.AndroidSdkPath = os.Getenv(Env.AndroidSdkEnv)

		fmt.Println(DEBUG_TAG, Env.AndroidSdkPath)

		// 从  build-tools 获取版本号， 如果没有配置版本号，取最大的版本号
		version, toolPath, err := findLatestSdkPath(path.Join(Env.AndroidSdkPath, "build-tools"))
		if err != nil {
			return err
		}
		Env.AndroidBuildToolVersion = version
		Env.AndroidBuildToolPath = toolPath
	}
	fmt.Println(Env)

	if nil != err {
		return err
	}
	return nil
}

func findLatestSdkPath(buildToolsPath string) (string, string, error) {
	files, err := ioutil.ReadDir(buildToolsPath)
	if err != nil {
		return "", "", err
	}

	latestVersion := ""
	for _, file := range files {
		if file.IsDir() {
			version := file.Name()
			if compareVersions(version, latestVersion) > 0 {
				latestVersion = version
			}
		}
	}

	if latestVersion == "" {
		return "", "", errors.New("no build tools found")
	}

	return latestVersion, filepath.Join(buildToolsPath, latestVersion), nil
}

func compareVersions(version1, version2 string) int {
	// 这里需要实现版本号比较的逻辑，可以使用 semver 或者其他方法
	return strings.Compare(version1, version2)
}
