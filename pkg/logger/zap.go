package logger

import (
	"errors"
	"os"
	"scaffold/pkg/config"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg *config.LogConfig) (err error) {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)

	// 创建编码器
	encoder := getEncoder()

	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return errors.New("l.UnmarshalText([]byte(cfg.Level)) failed")
	}

	var core zapcore.Core
	if cfg.Mode == "dev" {
		// 创建一个自定义的开发模式编码器配置
		devEncoderConfig := zap.NewDevelopmentEncoderConfig()
		// 应用相同的自定义时间格式
		devEncoderConfig.EncodeTime = customTimeEncoder
		devEncoderConfig.TimeKey = "time"
		devEncoderConfig.CallerKey = "caller"

		// 使用自定义配置创建控制台编码器
		consoleEncoder := zapcore.NewConsoleEncoder(devEncoderConfig)
		//开发模式,日志输出到终端
		core = zapcore.NewTee(
			//往终端写
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg := zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 - 15:04:05"))
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}

	// 添加文件写入器
	writers := []zapcore.WriteSyncer{zapcore.AddSync(lumberJackLogger)}

	writers = append(writers, zapcore.AddSync(os.Stdout))

	return zapcore.NewMultiWriteSyncer(writers...)
}
