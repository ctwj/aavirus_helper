package config

import (
	"context"
	"testing"
)

func TestLoad(t *testing.T) {
	s1 := `D:/work/gowork/src/register_tool/conf/config.yaml`
	//s1 := `D:/workspace/gopath/src/register_tool/conf/config.yaml`
	err := Load(context.Background(), s1)
	t.Log(App, err)
	// t.Log(Setting)
	// t.Log(Remote)
}
