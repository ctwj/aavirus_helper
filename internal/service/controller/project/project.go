package project

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/ctwj/aavirus_helper/internal/lib"
	"github.com/ctwj/aavirus_helper/internal/pkg/command"
	"github.com/ctwj/aavirus_helper/internal/pkg/config"
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

// 获取 apk 基本信息
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

// 反编译 apk
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
func (p *Project) BatchPack(apkdir string, list []string, mode string) interface{} {
	var result []string
	ctx := ctlCtx.Get()
	apkTool := path.Join(apkdir, "apktool.yml")

	log.Println(config.DEBUG_TAG, mode)

	for _, item := range list {
		if item != apkdir && item != apkTool {
			result = append(result, item)
		}
	}

	// 简单模式
	if mode == "single" {
		for i, removePath := range result {
			// 打包
			progress := fmt.Sprintf("%d/%d: %s ", i+1, len(result), path.Base(removePath))
			runtime.EventsEmit(*ctx, "progress", progress)
			command.NewCommand().DoPackAfterRemoveItem(apkdir, removePath)
		}
	}

	if mode == "group" {
		progress := "1/1: group"
		runtime.EventsEmit(*ctx, "progress", progress)
		command.NewCommand().DoPackAfterRemoveItems(apkdir, result)
	}

	runtime.EventsEmit(*ctx, "progress", "")

	return ""
}

// 批量打包 删除文件后进行打包， 需要排除掉 根目录，排除掉 apktool.yml
func (p *Project) BatchPermissionPack(apkdir string, list []string, mode string) interface{} {
	ctx := ctlCtx.Get()

	var files []interface{}
	var totalFile int

	// 简单模式
	if mode == "single" {
		for i, removePermission := range list {
			// 打包
			progress := fmt.Sprintf("%d/%d: %s ", i+1, len(list), lib.PermissionLastWord(removePermission))
			runtime.EventsEmit(*ctx, "progress", progress)
			dest, _ := command.NewCommand().DoPackAfterRemovePermission(apkdir, []string{removePermission})
			files = append(files, map[string]interface{}{"dest": dest, "remove": removePermission})
			totalFile = totalFile + 1
		}
	}

	if mode == "group" {
		progress := "1/1: group"
		runtime.EventsEmit(*ctx, "progress", progress)
		dest, _ := command.NewCommand().DoPackAfterRemovePermission(apkdir, list)
		files = append(files, map[string]interface{}{"dest": dest, "remove": strings.Join(list, ",")})
		totalFile = totalFile + 1
	}

	if mode == "cross" {
		crossList := lib.GenerateCrossCombinationData(list)
		for i, removePermissions := range crossList {
			// 打包
			name := lib.GeneratePackNameWithDelPermissions(removePermissions)
			progress := fmt.Sprintf("simple: %s %d/%d", name, i+1, len(crossList))
			runtime.EventsEmit(*ctx, "progress", progress)
			dest, _ := command.NewCommand().DoPackAfterRemovePermission(apkdir, removePermissions)
			files = append(files, map[string]interface{}{"dest": dest, "remove": strings.Join(list, ",")})
			totalFile = totalFile + 1
		}
	}

	htmlContent := lib.GenerateIndexHTML(mode, files)
	lib.WriteFile(path.Join(config.OutputDir, "index.html"), htmlContent)

	runtime.EventsEmit(*ctx, "progress", "")

	return ""
}

// 关闭打开的apk
func (p *Project) CloseApp(codedir string) interface{} {
	// 移除掉反编译的文件夹
	lib.RemovePaths([]string{codedir})
	return ""
}

// 获取 AndroidManifest.xml 的信息
func (p *Project) AndroidManifestInfo(codeDir string) interface{} {
	xmlPath := path.Join(codeDir, "AndroidManifest.xml")
	parse := lib.ParseAndroidManifest(xmlPath)
	content, _ := lib.ReadFileContent(xmlPath)
	return map[string]interface{}{
		"content": content,
		"parse":   parse,
	}
}

func (p *Project) OpenOutput() interface{} {
	command.NewCommand().OpenOutput()
	return ""
}
