package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/ctwj/aavirus_helper/internal/lib"
	"github.com/ctwj/aavirus_helper/internal/pkg/config"
	"github.com/ctwj/aavirus_helper/internal/pkg/tools"
)

type Command struct {
	Tool *tools.Tools
}

func NewCommand() *Command {
	return &Command{
		Tool: tools.NewTools(),
	}
}

func (c *Command) Run(cmds []string) ([]string, error) {
	var output []string

	for _, cmdStr := range cmds {
		lib.SendCommand2Front(fmt.Sprintf("#CMD: %v", cmdStr))

		// 输出内容
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", cmdStr)
		} else {
			cmd = exec.Command("bash", "-c", cmdStr)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("Error creating StdoutPipe for Cmd", err)
			lib.SendOutput2Front("Error creating StdoutPipe for Cmd" + err.Error())
			return output, err
		}

		if err := cmd.Start(); err != nil {
			log.Println("Error starting Cmd", err)
			lib.SendOutput2Front("Error starting Cmd:" + err.Error())
			return output, err
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			message := scanner.Text()

			log.Println(message)
			lib.SendOutput2Front(message)
			output = append(output, message)
			// 输出
		}

		if err := cmd.Wait(); err != nil {
			log.Println("Command failed:", err)
			lib.SendOutput2Front("Command failed:" + err.Error())
			return output, err
		}
	}

	// finish

	return output, nil
}

func (c *Command) OpenFolderCommand(path string) string {

	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = fmt.Sprintf("open \"%s\"", path)
	case "linux":
		cmd = fmt.Sprintf("xdg-open \"%s\"", path)
	case "windows":
		cmd = fmt.Sprintf("explorer \"%s\"", path)
	}
	return cmd
}

func (c *Command) CopyCommand(source, dest string) string {
	var cmd string
	if runtime.GOOS != "windows" {
		cmd = fmt.Sprintf("cp -r %s %s", source, dest)
	} else {
		cmd = fmt.Sprintf("xcopy %s %s /E", source, dest)
	}
	return cmd
}

func (c *Command) RemovePathCommand(path string) string {
	var cmd string
	if runtime.GOOS != "windows" {
		cmd = fmt.Sprintf("rm -rf %s", path)
	} else {
		cmd = fmt.Sprintf("rd /s /q %s", path)
	}
	return cmd
}

func (c *Command) PackCommand(codePath, destFile string) string {
	var cmd string
	java := c.Tool.Java
	apkTool := c.Tool.ApkTool

	cmd = fmt.Sprintf("%s -jar %s b %s -o %s", java, apkTool, codePath, destFile)

	return cmd
}

func (c *Command) EmptyFrameworkCommand() string {
	var cmd string
	java := c.Tool.Java
	apktool := c.Tool.ApkTool
	cmd = fmt.Sprintf("%s -jar  %s empty-framework-dir --force", java, apktool)
	return cmd
}

func (c *Command) CertGenerateCommand(randomPass, randomDName, jksPath string) string {
	var cmd string
	keytool := c.Tool.KeyTool

	cmd = fmt.Sprintf("%s -genkey -v -keystore  %s -alias my-alias -keyalg RSA -keysize 2048 -validity 10000 -storepass %s -keypass %s -dname \"%s\" -noprompt",
		keytool, jksPath, randomPass, randomPass, randomDName,
	)
	return cmd
}

func (c *Command) ApksignerCommand(apkPath, destPath, jksPath, randomPass string) string {
	var cmd string
	apksigner := c.Tool.ApkSigner

	cmd = fmt.Sprintf("%s sign --ks %s --ks-key-alias my-alias --ks-pass \"pass:%s\" --key-pass \"pass:%s\" --in %s --out %s",
		apksigner, jksPath, randomPass, randomPass, apkPath, destPath,
	)

	return cmd
}

func (c *Command) ZipalignCommand(sourceApk, destApk string) string {
	var cmd string
	zipalign := c.Tool.Zipalign

	cmd = fmt.Sprintf("%s -v 4 %s %s",
		zipalign, sourceApk, destApk)

	return cmd
}

func (c *Command) ApkInfoCommand(apkPath string) string {
	java := c.Tool.Java
	apkInfo := c.Tool.ApkInfo
	cmd := fmt.Sprintf("%s -jar %s %s", java, apkInfo, apkPath)
	return cmd
}

func (c *Command) DisassemblyCommand(apkPath string) (string, string) {
	java := c.Tool.Java
	apkTool := c.Tool.ApkTool
	fileName, _ := lib.GetFileNameWithoutExtension(apkPath)
	outputDir := path.Join(config.WorkDir, fileName)

	// 对输出文件夹 outputDir 进行检测，如果存在先清除
	if _, err := os.Stat(outputDir); err == nil {
		if err := os.RemoveAll(outputDir); err != nil {
			log.Printf("failed to remove existing output directory: %v", err)
			return outputDir, ""
		}
	}

	cmd := fmt.Sprintf("%s -jar %s d %s -o %s", java, apkTool, apkPath, outputDir)
	return outputDir, cmd
}

// java -jar ~/tools/apktool.jar b app-debug -o new-app-debug.apk

func (c *Command) DoDisassembly(apkPath string) (string, error) {
	outdir, cmd := c.DisassemblyCommand(apkPath)
	c.Run([]string{cmd})
	fmt.Println(cmd)
	return outdir, nil
}

func (c *Command) GetApkInfo(apkPath string) ([]string, error) {
	cmd := c.ApkInfoCommand(apkPath)
	result, err := c.Run([]string{cmd})
	if err != nil {
		return []string{}, err
	}
	fmt.Println(cmd)
	return result, nil
}

// 删除 removeItem 后重新打包
func (c *Command) DoPackAfterRemoveItem(codePath string, removeItem string) (string, error) {

	// 第一步， 复制代码到 tmp 目录
	tmp := lib.GenerateRandomString(8)
	targetCodeDir := path.Join(config.TmpDir, tmp)

	// 第二步， 移除目录
	removeDir := strings.Replace(removeItem, codePath, targetCodeDir, -1)
	c.Run([]string{
		c.CopyCommand(codePath, targetCodeDir),
		c.RemovePathCommand(removeDir),
	})

	// 第三步， 随机包名
	randomPackName := lib.GenerateRandomPackName()
	lib.ChangePackName(targetCodeDir, randomPackName)

	// 第四部， 打包, 签名
	relativePath := strings.Replace(removeItem, codePath, "", -1)
	desFileName := lib.GenerateTargetFileName(relativePath, removeItem)

	destFile_1 := path.Join(config.TmpDir, desFileName+"_1.apk")

	// 签名相关
	randomPass := lib.GenerateRandomString(16)
	jksPath := path.Join(config.TmpDir, randomPass+".jks")
	randomDName := lib.GenerateRandomDName()

	// 打包签名后的文件
	// destFile_2 := path.Join(config.TmpDir, desFileName+"_2.apk")
	destFile_3 := path.Join(config.OutputDir, desFileName+".apk")

	c.Run([]string{
		c.EmptyFrameworkCommand(),
		c.PackCommand(targetCodeDir, destFile_1),
		c.CertGenerateCommand(randomPass, randomDName, jksPath),
		c.ApksignerCommand(destFile_1, destFile_3, jksPath, randomPass),
		// c.ZipalignCommand(destFile_2, destFile_3), // 不对齐
	})

	lib.RemovePaths([]string{
		targetCodeDir,
		jksPath,
		destFile_1,
	})

	return destFile_3, nil
}

func (c *Command) OpenOutput() {
	cmd := c.OpenFolderCommand(config.OutputDir)
	c.Run([]string{cmd})
}
