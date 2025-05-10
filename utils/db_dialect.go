package utils

import "gorm.io/gorm"

func ResolveDBRandomFunc(orm *gorm.DB) string {
	dialect := orm.Dialector.Name()
	// 根据不同数据库使用对应的随机函数
	var randomFunc string
	switch dialect {
	case "mysql":
		randomFunc = "RAND()"
	case "postgres", "postgresql":
		randomFunc = "RANDOM()"
	case "sqlite":
		randomFunc = "RANDOM()"
	default:
		randomFunc = "RAND()" // 默认使用MySQL语法
	}
	return randomFunc
}
