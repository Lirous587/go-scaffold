package cmd

import (
	"fmt"
	"scaffold/pkg/config"
	"scaffold/pkg/engine"
	"scaffold/pkg/logger"
	"scaffold/pkg/repository/redis"

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

	if err = redis.Init(&config.Cfg.Redis); err != nil {
		return errors.WithMessage(err, "redis模块初始化失败")
	}

	return nil
}

func Main() {
	if err := setting(); err != nil {
		fmt.Printf("%+v", err)
		return
	}

	router := engine.Init(&config.Cfg.App)

	engine.Run(router, config.Cfg.App.Port)
}
