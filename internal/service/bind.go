package service

import (
	"github.com/ctwj/aavirus_helper/internal/service/controller/project"
	"github.com/ctwj/aavirus_helper/internal/service/controller/task"
	"github.com/ctwj/aavirus_helper/internal/service/controller/upload"
)

func Bind() []interface{} {
	return []interface{}{
		&task.Task{},
		&project.Project{},
		&upload.Adb{},
		&upload.Web{},
	}
}
