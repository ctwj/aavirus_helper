package tools

import (
	"path"

	"github.com/ctwj/aavirus_helper/internal/pkg/config"
)

type Tools struct {
	Java      string
	ApkTool   string
	KeyTool   string
	ApkSigner string
	JarSigner string
	Zipalign  string
	ApkInfo   string
}

func NewTools() *Tools {

	env := config.Env

	return &Tools{
		Java:      "java",
		ApkTool:   path.Join(config.AppDir, env.ToolDir, env.ApkToolName),
		KeyTool:   "keytool",
		ApkSigner: path.Join(env.AndroidBuildToolPath, "apksigner"),
		JarSigner: "jarsigner",
		Zipalign:  path.Join(env.AndroidBuildToolPath, "zipalign"),
		ApkInfo:   path.Join(config.AppDir, env.ToolDir, "GetMoreAPKInfo.jar"),
	}
}
