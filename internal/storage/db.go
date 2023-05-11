package storage

import (
	"fmt"

	"github.com/dreaminglwj/sage/internal/resource"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	xormLog "xorm.io/xorm/log"
	"xorm.io/xorm/names"

	"github.com/dreaminglwj/sage/internal/conf"
	"github.com/dreaminglwj/sage/internal/plugin/log"
)

func NewEngine(cfg *conf.Config, logger *log.Logger) (*xorm.Engine, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Name,
	)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if cfg.App.Debug {
		logger.SetLevel(xormLog.LOG_DEBUG)
	} else {
		logger.SetLevel(xormLog.LOG_WARNING)
	}
	logger.ShowSQL(true)
	engine.SetLogger(logger)
	engine.SetMaxOpenConns(20)
	engine.SetMapper(names.GonicMapper{})
	resource.Register(engine)
	return engine, nil
}
