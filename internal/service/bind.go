package service

import (
	"github.com/ctwj/aavirus_helper/internal/service/controller/task"
)

func Bind() []interface{} {
	return []interface{}{
		&task.Task{},
	}
}