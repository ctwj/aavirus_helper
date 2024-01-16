package command

import (
	"bufio"
	"fmt"
	"log"
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

func (c *Command) Run(cmds []string) error {
	for i, cmdStr := range cmds {
		println(i)

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
			return err
		}

		if err := cmd.Start(); err != nil {
			log.Println("Error starting Cmd", err)
			return err
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			message := scanner.Text()

			log.Println(message)
			// 输出
		}

		if err := cmd.Wait(); err != nil {
			log.Println("Command failed:", err)
			return err
		}
	}

	// finish

	return nil
}

func (c *Command) DisassemblyCommand(apkPath string) string {
	java := c.Tool.Java
	apkTool := c.Tool.ApkTool
	fileName, _ := lib.GetFileNameWithoutExtension(apkPath)
	outputDir := path.Join(config.WorkDir, fileName)
	cmd := fmt.Sprintf("%s -jar %s d %s -o %s", java, apkTool, apkPath, outputDir)
	return cmd
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

func (c *Command) DoDisassembly(apkPath string) error {
	cmd := c.DisassemblyCommand(apkPath)
	c.Run([]string{cmd})
	fmt.Println(cmd)
	return nil
}
