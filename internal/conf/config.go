package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

// NewConfig config 初始化
func NewConfig(filepath string) (*Config, error) {
	c := config.New(
		config.WithSource(
			env.NewSource(""),
			file.NewSource(filepath),
		),
	)
	defer func() {
		if err := c.Close(); err != nil {
			panic(err)
		}
	}()

	if err := c.Load(); err != nil {
		return nil, err
	}

	var bc Config
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// Config 全局Config
type Config struct {
	App    *App
	Root   *RootUser
	Server *struct {
		GRPC *server
	}
	Database *Database
	Log      *Log
	Jaeger   *Jaeger
	DBTables *DBTables

	OutFileName string
}

// App 应用基本配置
type App struct {
	Name      string
	Version   string
	Level     uint
	Component string
	Debug     bool
	Env       Env
}
type Env string

const (
	EnvProduction Env = "prod"
	EnvTest       Env = "test"
	EnvDev        Env = "dev"
	EnvLocal      Env = "local"
)

func (a *App) IsProduction() bool {
	return a.Env == EnvProduction
}

type RootUser struct {
	Account  string
	Email    string
	Password string
}

type server struct {
	Network string
	Addr    string
	Timeout int
}

// Database 数据库配置
type Database struct {
	Host     string
	User     string
	Password string
	Name     string
}

// Log  日志配置
type Log struct {
	Type  string
	Level string
}

type Jaeger struct {
	Endpoint string
}

type DBTables struct {
	Tables map[string][]string
}
