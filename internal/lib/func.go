package lib

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ctlCtx "github.com/ctwj/aavirus_helper/internal/service/context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ==================================================
// 生成随机字符串
func GenerateRandomString(len int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, len)
	for i := range b {
		b[i] = charset[rand.Intn(len)]
	}
	return string(b)
}

// ==================================================
// 生成随机dname
// 示例 CN=Your Name, OU=Your Organization, O=Your Company, L=Your City, ST=Your State, C=Your Country
func GenerateRandomDName() string {
	return fmt.Sprintf("CN=%s, OU=%s, O=%s, L=%s, ST=%s C=%s",
		GenerateRandomString(2),
		GenerateRandomString(4),
		GenerateRandomString(6),
		GenerateRandomString(6),
		GenerateRandomString(6),
		GenerateRandomString(6),
	)
}

// ==================================================
// 生成随机包名
func GenerateRandomPackName() string {
	return fmt.Sprintf("%s.%s.%s", GenerateRandomString(6), GenerateRandomString(6), GenerateRandomString(6))
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

// ChangePackName 修改apk包名
func ChangePackName(codePath, newPackName string) {
	manifestFile := path.Join(codePath, "AndroidManifest.xml")
	apktoolFile := path.Join(codePath, "apktool.yml")

	// 读取文件内容
	content, err := os.ReadFile(manifestFile)
	if err != nil {
		fmt.Println("Error reading AndroidManifest.xml:", err)
		return
	}

	// step 1, change  AndroidManifest.xml 中包名
	// 提取旧包名
	oldPackNamePattern := regexp.MustCompile(`package="(.*?)"`)
	oldPackNameMatches := oldPackNamePattern.FindStringSubmatch(string(content))
	if len(oldPackNameMatches) < 2 {
		fmt.Println("Old package name not found in AndroidManifest.xml")
		return
	}

	oldPackName := oldPackNameMatches[1]
	// 替换包名
	newContent := regexp.MustCompile(oldPackName).ReplaceAllString(string(content), newPackName)
	// 将修改后的内容写回文件
	if err := os.WriteFile(manifestFile, []byte(newContent), os.ModePerm); err != nil {
		fmt.Println("Error writing AndroidManifest.xml:", err)
		return
	}

	// 读取文件内容
	configContent, err := os.ReadFile(apktoolFile)
	if err != nil {
		fmt.Println("Error reading AndroidManifest.xml:", err)
		return
	}

	// 替换内容
	oldStr := "renameManifestPackage: null"
	newStr := fmt.Sprintf("renameManifestPackage: %s", newPackName)
	newConfigContent := strings.ReplaceAll(string(configContent), oldStr, newStr)

	// 将修改后的内容写回文件
	if err := os.WriteFile(apktoolFile, []byte(newConfigContent), os.ModePerm); err != nil {
		fmt.Println("Error writing AndroidManifest.xml:", err)
		return
	}

	// step 1，change  apktool.yml 中包名
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func GenerateTargetFileName(relativePath, removePath string) string {
	isFile := IsFile(removePath)
	preSuff := "_d"
	if isFile {
		preSuff = "_f"
	}

	// 使用 strings.Split 将字符串切分为数组
	parts := strings.Split(relativePath, "/")

	// 移除第一个元素（如果为空字符串）
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}

	// 取出最后一个字符串
	lastPart := parts[len(parts)-1]

	// 将其他部分连接为一个字符串，使用 "-"
	var joinedParts string
	if len(parts) > 1 {
		joinedParts = strings.Join(parts[1:len(parts)-1], "-")
	} else {
		joinedParts = "d" + parts[0]
	}

	// 将连接的字符串与最后一个字符串连接起来
	result := joinedParts + preSuff + lastPart

	return result
}
