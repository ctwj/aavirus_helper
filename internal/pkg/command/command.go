package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"

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

	for i, cmdStr := range cmds {
		if len(cmds) > 1 {
			lib.SendCommand2Front(fmt.Sprintf("#  ==== Step: %v ====", i+1))
		}

		lib.SendCommand2Front(cmdStr)

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

/**
keytool -genkey -v -keystore my.jks -alias my-alias -keyalg RSA -keysize 2048 -validity 10000

# 生成证书
keytool -genkey -v -keystore my.jks -alias my-alias -keyalg RSA -keysize 2048 -validity 10000 -storepass random-keystore-password -keypass random-key-password -dname "CN=Your Name, OU=Your Organization, O=Your Company, L=Your City, ST=Your State, C=Your Country" -noprompt

# 签名
~/Android/Sdk/build-tools/34.0.0/apksigner sign --ks my.jks --ks-key-alias my-alias --ks-pass "pass:random-keystore-password" --key-pass "pass:random-key-password" --in new-app-debug.apk --out result.apk

# 验签
jarsigner -verify -verbose -certs result.apk

~/Android/Sdk/build-tools/34.0.0/zipalign -v 4 new-app-debug.apk aligned.apk
**/

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
