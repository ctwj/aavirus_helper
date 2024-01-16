package task

import (
	"path"

	"github.com/ctwj/aavirus_helper/internal/pkg/command"
	"github.com/ctwj/aavirus_helper/internal/pkg/config"
)

type Task struct {
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Create() {
	command.NewCommand().DoDisassembly(path.Join(config.AppDir, "tmp", "realapk_v1.0.6_prod_01-15_00-45.apk"))
}
