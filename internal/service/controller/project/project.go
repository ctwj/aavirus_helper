package project

import (
	"path"

	"github.com/ctwj/aavirus_helper/internal/lib"
	"github.com/ctwj/aavirus_helper/internal/pkg/command"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	ctlCtx "github.com/ctwj/aavirus_helper/internal/service/context"
)

type Project struct {
}

func NewProject() *Project {
	return &Project{}
}

// OpenFile 选择一个文件 返回 {file: "xxx.apk"}
func (p *Project) OpenFile() (interface{}, error) {
	ctx := ctlCtx.Get()
	file, err := runtime.OpenFileDialog(*ctx, runtime.OpenDialogOptions{
		Title: "Open an apk file",
		Filters: []runtime.FileFilter{{
			DisplayName: "Apk Files (*.apk)",
			Pattern:     "*.apk",
		}},
	})

	if err != nil {
		return "", err
	}
	if file == "" {
		return "{\"file\":\"\"}", nil
	}
	// result, _ := json.Marshal(map[string]string{"file": file})
	return map[string]string{"file": file}, nil
}

// 显示反编译后的文件列表
func (p *Project) FileList(dir string) interface{} {
	// 获取树状结构的文件列表， 包含文件名， 文件大小， 文件修改时间，
	// 如果是是文件夹，还需要知道文件夹中文件的总大小，和总文件个数
	list, _ := lib.FileList(dir) // 获取文件列表
	lib.CalculateDirSize(&list)  // 计算文件夹大小和文件个数
	return list
}

func (p *Project) GetApkInfo(apkPath string) interface{} {
	info, err := command.NewCommand().GetApkInfo(apkPath)
	if err != nil {
		return map[string]interface{}{
			"info": "",
			"err":  err.Error(),
		}
	}
	return map[string]interface{}{
		"info": info,
		"err":  nil,
	}
}

func (p *Project) Disassemble(apkPath string) interface{} {
	outdir, err := command.NewCommand().DoDisassembly(apkPath)
	if err != nil {
		return map[string]interface{}{
			"info":   "",
			"outdir": outdir,
			"err":    err.Error(),
		}
	}
	return map[string]interface{}{
		"info":   "success",
		"outdir": outdir,
		"err":    nil,
	}
}

// 批量打包 删除文件后进行打包， 需要排除掉 根目录，排除掉 apktool.yml
func (p *Project) BatchPack(apkdir string, list []string) interface{} {
	var result []string
	apkTool := path.Join(apkdir, "apktool.yml")

	for _, item := range list {
		if item != apkdir && item != apkTool {
			result = append(result, item)
		}
	}

	for _, removePath := range result {
		// 打包
		command.NewCommand().DoPackAfterRemoveItem(apkdir, removePath)
	}

	return ""
}
