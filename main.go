package main

//func setting() error {
//	var err error
//	if err = config.Init(); err != nil {
//		return errors.WithMessage(err, "config模块初始化失败")
//	}
//
//	if err = logger.Init(&config.Cfg.Log); err != nil {
//		return errors.WithMessage(err, "config模块初始化失败")
//	}
//
//	// 初始化验证器
//	if err = validator.Init(); err != nil {
//		return errors.WithMessage(err, "validator模块初始化失败")
//	}
//
//	if err = repository.RedisInit(&config.Cfg.Redis); err != nil {
//		return errors.WithMessage(err, "redis模块初始化失败")
//	}
//
//	if err = repository.GormInit(&config.Cfg.DB); err != nil {
//		return errors.WithMessage(err, "db模块初始化失败")
//	}
//
//	if err = email.Init(&config.Cfg.Email); err != nil {
//		return errors.WithMessage(err, "email模块初始化失败")
//	}
//
//	return nil
//}
//
//func main() {
//	if err := setting(); err != nil {
//		log.Fatalf("%+v", err)
//		return
//	}
//
//}
