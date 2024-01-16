package config

const (
	EnvAppDirName = "APP_DIR"

	DEBUG_TAG = "\033[0;33m[DEBUG]\033[0m"
	ERROR_TAG = "\033[0;31m[ERROR]\033[0m"
	INFO_TAG  = "\033[0;34m[INFO]\033[0m"
)

// 程序根目录
var AppDir string

var TmpDir string

var WorkDir string

var OutputDir string
