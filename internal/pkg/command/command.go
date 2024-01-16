package command

import "github.com/ctwj/aavirus_helper/internal/pkg/tools"

type Command struct {
	Tool *tools.Tools
}

func NewCommand() *Command {
	return &Command{
		Tool: tools.NewTools(),
	}
}

func (c *Command) Disassembly(apkPath string, outDir string) error {
	return nil
}
