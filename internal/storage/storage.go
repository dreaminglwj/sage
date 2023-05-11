package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	dockertest "github.com/ory/dockertest"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/dreaminglwj/sage/internal/resource"

	"github.com/dreaminglwj/sage/internal/conf"
	log2 "github.com/dreaminglwj/sage/internal/plugin/log"
)

type Storage struct {
	engine *xorm.Engine
	logger *log.Helper
}

// NewStorage is the constructor of Storage
func NewStorage(engine *xorm.Engine, logger *log.Helper) *Storage {
	return &Storage{
		engine: engine,
		logger: logger,
	}
}

func (s *Storage) DB(ctx context.Context) *xorm.Session {
	return s.engine.Context(ctx)
}

func (s *Storage) TX(ctx context.Context, fn func(db *xorm.Session) error) error {
	session := s.engine.NewSession().Context(ctx)
	defer func() {
		if err := session.Close(); err != nil {
			s.logger.WithContext(ctx).Error("xorm close transaction session error: ", err)
		}
	}()
	if err := session.Begin(); err != nil {
		return err
	}
	if err := fn(session); err != nil {
		return err
	}
	return session.Commit()
}

func (s *Storage) AutoMigrate() error {
	//return s.engine.Sync2(
	//
	//)
	return nil
}

func (s *Storage) addAdmin() error {
	// now := time.Now()
	// admin := &model.Manager{
	// 	Id:        "82f8488e-a2f0-cf93-5410-078f404fd2b9",
	// 	Account:   "admin",
	// 	Password:  "$2a$10$tm6xtWJe.wEJ/egbZIm3VOX6Cl8cpnGJESnI0tCPZ09Zwj90Mr40W",
	// 	Name:      "admin",
	// 	Email:     "admin@wst.com",
	// 	Phone:     "17777777777",
	// 	Avatar:    "",
	// 	IsValid:   true,
	// 	CreatedAt: now,
	// 	UpdatedAt: now,
	// }

	db := s.DB(context.TODO())
	_, err := db.Exec("INSERT INTO `manager`(`id`,`account`,`password`,`name`,`email`,`phone`,`avatar`,`is_valid`,`created_at`,`updated_at`) VALUES('82f8488e-a2f0-cf93-5410-078f404fd2b9','admin', '$2a$10$tm6xtWJe.wEJ/egbZIm3VOX6Cl8cpnGJESnI0tCPZ09Zwj90Mr40W', 'admin','admin@wst.com', '17777777777', '', 1, '2023-04-08 12:00:00', '2023-04-08 12:00:00')")
	if err != nil {
		return err
	}
	return nil
}

const createTestDatabaseSql = "CREATE DATABASE IF NOT EXISTS docker DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;"

// 初始化 Docker mysql 容器
func innerDockerMysql(img, version string, log *log.Helper) (*conf.Database, closer) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %+v testStorage", err)
	}

	// pulls an image, creates a container based on it and runs it
	dbResource, err := pool.Run(img, version, []string{"MYSQL_ROOT_PASSWORD=secret", "MYSQL_ROOT_HOST=%"})
	if err != nil {
		log.Fatalf("Could not start dbResource: %+v testStorage", err)
	}
	hostPort := dbResource.GetHostPort("3306/tcp")
	config := conf.Database{
		Host:     hostPort,
		User:     "root",
		Password: "secret",
		Name:     "docker",
	}
	conStr := fmt.Sprintf("%s:%s@tcp(%s)/mysql?charset=utf8mb4&parseTime=true", config.User, config.Password, config.Host)
	if err := pool.Retry(func() error {
		db, err := sql.Open("mysql", conStr)
		if err != nil {
			return err
		}
		if err := db.Ping(); err != nil {
			return err
		}
		_, err = db.Exec(createTestDatabaseSql)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %+v testStorage", err)
	}

	// 回调函数关闭容器
	return &config, closer{fn: func() error {
		if err := pool.Purge(dbResource); err != nil {
			log.Fatalf("Could not purge dbResource: %+v testStorage", err)
			return err
		}
		return nil
	}}
}

type closer struct {
	fn func() error
}

func (c closer) Close() error {
	return c.fn()
}

func NewTestStorage() (*Storage, error) {
	logger, err := log2.NewLogger(&conf.Config{App: &conf.App{
		Env: "test",
	}})
	if err != nil {
		panic(err)
	}
	helper := log.NewHelper(logger)
	//这个镜像支持MAC M1处理器
	cfg, cls := innerDockerMysql("mysql", "oracle", helper)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Name,
	)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true)
	engine.SetMapper(names.GonicMapper{})
	resource.Register(engine)
	resource.Register(cls)
	s := NewStorage(engine, helper)
	if err := s.AutoMigrate(); err != nil {
		return nil, err
	}
	err = s.addAdmin()
	if err != nil {
		return nil, err
	}
	return s, nil
}
