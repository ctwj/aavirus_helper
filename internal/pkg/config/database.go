package config

import (
	"context"
	"time"

	"github.com/spf13/viper"
)

const (
	sqlTypeMysql = "mysql"
)

// 数据库引擎
type SqlEngine string

const (
	SqlEngineDefault SqlEngine = `default`
)

type SqlItem struct {
	Enable          bool `mapstructure:"enable"`
	Label           string
	Type            string        `mapstructure:"type"`
	TableOptions    string        `mapstructure:"table_opts"`
	Master          string        `mapstructure:"master"`
	Slaves          []string      `mapstructure:"slaves"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Passwd          string        `mapstructure:"password"`
	Schema          string        `mapstructure:"schema"`
	SSLMode         bool          `mapstructure:"ssl_mode"`
	Charset         string        `mapstructure:"charset"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnectTimeout  time.Duration `mapstructure:"conn_time_out"`
	ConnectRetries  int           `mapstructure:"conn_retries"`
	ConnectBackoff  time.Duration `mapstructure:"conn_retry_backoff"`
	LogSQL          bool          `mapstructure:"log_sql"`
	LogPath         string        `mapstructure:"log_path"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_life_time"`
}

var (
	SqlDB []*SqlItem
)

func LoadSqlDBConfig(ctx context.Context, workDir string) error {
	sqlCfg := make(map[string]*SqlItem)
	err := viper.Sub("database").Unmarshal(&sqlCfg)
	if err != nil {
		return err
	}
	for key, item := range sqlCfg {
		item.Label = key
		item.ConnMaxLifetime = item.ConnMaxLifetime * time.Second
		if item.Type == sqlTypeMysql {
			if item.ConnMaxLifetime == 0 {
				item.ConnMaxLifetime = 3 * time.Second
			}
		}

		item.ConnectBackoff *= time.Second //3 * time.Second

		SqlDB = append(SqlDB, item)
	}
	return nil
}
