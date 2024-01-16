package project

import "github.com/ctwj/aavirus_helper/internal/lib"

type Project struct {
}

func NewProject() *Project {
	return &Project{}
}

// 显示反编译后的文件列表
func (p *Project) FileList(apkPath string, dir string) interface{} {
	// 获取树状结构的文件列表， 包含文件名， 文件大小， 文件修改时间，
	// 如果是是文件夹，还需要知道文件夹中文件的总大小，和总文件个数
	list, _ := lib.FileList(apkPath, dir)
	return list
}
