package email

import (
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

type mailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	CC       string
}

type mailer struct {
	dialer *gomail.Dialer
}

var instance Mailer

var config mailerConfig

func Init() {
	UpdateConfig()
	validateEnv()
}

func validateEnv() {
	// 校验必填项
	if config.Host == "" ||
		config.Port == 0 ||
		config.Username == "" ||
		config.Password == "" ||
		config.From == "" {
		panic(errors.New("email config: 环境变量缺失，必填项不能为空"))
	}
}

func UpdateConfig() {
	portStr := os.Getenv("EMAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	config = mailerConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     port,
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
		FromName: os.Getenv("EMAIL_FROM_NAME"),
		CC:       os.Getenv("EMAIL_CC"),
	}

	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	instance = &mailer{
		dialer: dialer,
	}
}
