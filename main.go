package main

import (
	"fmt"
	"scaffold/internal/cmd"
	"scaffold/pkg/config"
	"scaffold/pkg/i18n"
	"scaffold/pkg/logger"
	"scaffold/pkg/repository/db"
	"scaffold/pkg/repository/redis"
	"scaffold/pkg/validator"

	"github.com/pkg/errors"
)

func setting() error {
	var err error
	if err = config.Init(); err != nil {
		return errors.WithMessage(err, "config模块初始化失败")
	}

	if err = logger.Init(&config.Cfg.Log); err != nil {
		return errors.WithMessage(err, "config模块初始化失败")
	}

	// 初始化i18n (确保在validator之前)
	i18n.Init()

	// 初始化验证器
	if err = validator.Init(); err != nil {
		return errors.WithMessage(err, "validator模块初始化失败")
	}

	if err = redis.Init(&config.Cfg.Redis); err != nil {
		return errors.WithMessage(err, "redis模块初始化失败")
	}

	if err = db.Init(&config.Cfg.DB); err != nil {
		return errors.WithMessage(err, "db模块初始化失败")
	}

	return nil
}

func main() {
	if err := setting(); err != nil {
		fmt.Printf("%+v", err)
		return
	}

	cmd.Main()
}
