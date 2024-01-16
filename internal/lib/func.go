package lib

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// ==================================================
// 生成随机字符串
func GenerateRandomString(len int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, len)
	for i := range b {
		b[i] = charset[rand.Intn(len)]
	}
	return string(b)
}

// ==================================================
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

// ==================================================
// 如果目录不存在创建目录
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

func GetFileNameWithoutExtension(filePath string) (string, error) {
	// 获取文件名（带后缀）
	fileName := filepath.Base(filePath)

	// 去掉文件后缀
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

	return fileNameWithoutExt, nil
}

// ==================================================
// FileList 获取树状结构的文件列表
type FileInfo struct {
	Name         string
	Size         int64
	IsDir        bool
	ModTime      time.Time
	TotalSize    int64
	TotalFileNum int
}

// FileList gets the tree-like structure of file list
func FileList(apkPath string, dir string) ([]FileInfo, error) {
	var fileList []FileInfo
	err := listFiles(apkPath, dir, &fileList)
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

// listFiles recursively traverses the directory
func listFiles(apkPath string, dir string, fileList *[]FileInfo) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var dirInfo FileInfo
	dirInfo.Name = filepath.Base(dir)
	dirInfo.IsDir = true
	dirInfo.ModTime = time.Now()

	for _, file := range files {
		fileInfo := FileInfo{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		}

		if file.IsDir() {
			// Recursively process subdirectories
			subDir := filepath.Join(dir, file.Name())
			err := listFiles(apkPath, subDir, fileList)
			if err != nil {
				return err
			}
		} else {
			fileInfo.Size, err = fileSize(file)
			if err != nil {
				return err
			}
		}

		dirInfo.TotalSize += fileInfo.Size
		dirInfo.TotalFileNum++

		*fileList = append(*fileList, fileInfo)
	}

	*fileList = append(*fileList, dirInfo)
	return nil
}

func fileSize(file fs.DirEntry) (int64, error) {
	info, err := file.Info()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
