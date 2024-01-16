package biz

import (
	"path"
	"path/filepath"

	"github.com/ctwj/aavirus_helper/internal/lib"
)

type Task struct {
	TaskId   string // 新任务id，需要一个唯一路径
	Path     string // apk 加载路径
	CodePath string // 反汇编代码路径
}

func CreateTask(apkPath string) *Task {

	fileName := path.Base(apkPath)
	fileNameWithoutExt := fileName[:len(fileName)-len(path.Ext(fileName))]

	return &Task{
		TaskId:   lib.GenerateRandomString(4),
		Path:     apkPath,
		CodePath: filepath.Join(".", "workspace", fileNameWithoutExt),
	}
}
