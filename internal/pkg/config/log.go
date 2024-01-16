package config

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	Log = struct {
		SplitSize  uint
		FilePath   string
		Level      string
		ToConsole  bool
		KeepDays   uint
		MaxBackups int
		Compress   bool
	}{}
)

func LoadLogConfig(ctx context.Context, dir string) error {
	logViper := viper.Sub("log")
	Log.SplitSize = logViper.GetUint("split_size")
	filename := logViper.GetString("file_path")
	if filename != "" {
		if filepath.IsAbs(filename) {
			Log.FilePath = filename
		} else {
			Log.FilePath = filepath.Join(dir, filename)
		}
	} else {
		Log.FilePath = filepath.Join(dir, fmt.Sprintf("%s.log", App.Name))
	}

	Log.Level = logViper.GetString("level")
	Log.ToConsole = logViper.GetBool("to_console")
	Log.Compress = logViper.GetBool("compress")
	Log.KeepDays = logViper.GetUint("keep_days")
	Log.MaxBackups = logViper.GetInt("max_backups")
	return nil
}
