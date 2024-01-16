package lib

import (
	"fmt"
	"math/rand"
	"os"
)

func GenerateRandomString(len int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, len)
	for i := range b {
		b[i] = charset[rand.Intn(len)]
	}
	return string(b)
}

// 从环境变量获取 android sdk 路径
func FindAndroidSDK() (string, error) {
	// 首先检查 ANDROID_HOME 环境变量
	androidHome := os.Getenv("ANDROID_HOME")
	if androidHome != "" {
		return androidHome, nil
	}

	// 如果 ANDROID_HOME 未设置，则检查 ANDROID_SDK_ROOT 环境变量
	androidSDKRoot := os.Getenv("ANDROID_SDK_ROOT")
	if androidSDKRoot != "" {
		return androidSDKRoot, nil
	}

	// 如果两者都未设置，则返回错误
	return "", fmt.Errorf("Android SDK not found. Set either ANDROID_HOME or ANDROID_SDK_ROOT environment variable.")
}

func CreateDirectoryIfNotExists(dirPath string) error {
	// 检查目录是否存在
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
		fmt.Printf("Created directory: %s\n", dirPath)
	} else if err != nil {
		return err
	}

	return nil
}
