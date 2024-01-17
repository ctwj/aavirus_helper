package lib

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	ctlCtx "github.com/ctwj/aavirus_helper/internal/service/context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
// FileInfo 获取文件或目录信息
type FileInfo struct {
	Label        string     `json:"label"`
	Value        string     `json:"value"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Size         int64      `json:"size"`
	IsDir        bool       `json:"isDir"`
	Path         string     `json:"path"`
	TotalSize    int64      `json:"totalSize,omitempty"`
	TotalFileNum int        `json:"totalFileNum,omitempty"`
	HumanSize    string     `json:"humanSize,omitempty"`
	Children     []FileInfo `json:"children,omitempty"`
}

func FileList(dir string) (FileInfo, error) {
	var dirInfo FileInfo
	info, err := os.Stat(dir)
	if err != nil {
		return dirInfo, err
	}
	if !info.IsDir() { // 该函数值遍历文件夹
		return dirInfo, nil //
	}
	dirInfo.Name = filepath.Base(dir)
	dirInfo.IsDir = true
	dirInfo.Path = dir
	dirInfo.Size = 0
	dirInfo.TotalSize = 0
	dirInfo.TotalFileNum = 0
	dirInfo.Label = dirInfo.Name
	dirInfo.Value = dirInfo.Path
	dirInfo.Key = dirInfo.Path

	files, err := os.ReadDir(dir)
	if err != nil {
		return dirInfo, err
	}
	var children []FileInfo
	for _, file := range files {

		filePath := filepath.Join(dir, file.Name())
		var item FileInfo
		if file.IsDir() {
			item, _ = FileList(filePath)
		} else {
			item.Path = filePath
			info, _ := os.Stat(filepath.Join(dir, file.Name()))
			item.Name = file.Name()
			item.Size = info.Size()
			item.Label = item.Name
			item.Value = item.Path
			item.Key = item.Path
			item.IsDir = false
		}
		children = append(children, item)
	}
	dirInfo.Children = children
	return dirInfo, nil
}

// 计算 FileList 中文件个数 和 目录大小
func CalculateDirSize(root *FileInfo) (int64, int) {
	var totalSize int64
	var totalFileNum int
	var size int64
	size = root.Size
	if root.IsDir {

		// 计算子目录的大小
		for i := range root.Children {
			if root.Children[i].IsDir {
				dirTotalSize, dirTotalFileNum := CalculateDirSize(&root.Children[i])
				totalSize = totalSize + dirTotalSize
				totalFileNum = totalFileNum + dirTotalFileNum
			} else {
				totalSize = totalSize + root.Children[i].Size
				totalFileNum = totalFileNum + 1
			}
		}
		size = totalSize
	}
	root.TotalSize = totalSize
	root.TotalFileNum = totalFileNum
	root.HumanSize = HumanFileSize(float64(size))
	return totalSize, totalFileNum
}

// ==================================================
// 发送命令到前端
func SendCommand2Front(cmd string) {
	ctx := ctlCtx.Get()
	runtime.EventsEmit(*ctx, "command", cmd)
}

// ==================================================
// 发送结果到前端
func SendOutput2Front(cmd string) {
	ctx := ctlCtx.Get()
	runtime.EventsEmit(*ctx, "message", cmd)
}

// ==================================================
// 格式化文件大小
var (
	suffixes [5]string
)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func HumanFileSize(size float64) string {
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}
