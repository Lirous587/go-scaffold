package logger

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapConfig struct {
	level      string
	fileName   string
	maxSize    int
	maxAge     int
	maxBackups int
}

var config zapConfig

func UpdateConfig() {
	level := os.Getenv("LOG_LEVEL")
	fileName := os.Getenv("LOG_FILENAME")
	maxSizeStr := os.Getenv("LOG_MAX_SIZE")
	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil {
		panic(err)
	}
	maxAgeStr := os.Getenv("LOG_MAX_AGE")
	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		panic(err)
	}
	maxBackupsStr := os.Getenv("LOG_MAX_BACKUPS")
	maxBackups, err := strconv.Atoi(maxBackupsStr)
	if err != nil {
		panic(err)
	}

	if level == "" || fileName == "" || maxSize == 0 || maxAge == 0 || maxBackups == 0 {
		panic(errors.New("zap配置文件加载失败"))
	}

	config = zapConfig{
		level:      level,
		fileName:   fileName,
		maxSize:    maxSize,
		maxAge:     maxAge,
		maxBackups: maxBackups,
	}

}

func Init() (err error) {
	_ = godotenv.Load()

	UpdateConfig()
	writeSyncer := getLogWriter()
	// 创建编码器
	encoder := getEncoder()

	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(config.level))
	if err != nil {
		return errors.New("l.UnmarshalText([]byte(cfg.Level)) failed")
	}

	var core zapcore.Core
	core = zapcore.NewCore(encoder, writeSyncer, l)

	lg := zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(lg)
	return
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.fileName,
		MaxSize:    config.maxSize,
		MaxBackups: config.maxBackups,
		MaxAge:     config.maxAge,
	}

	// 添加文件写入器
	writers := []zapcore.WriteSyncer{zapcore.AddSync(lumberJackLogger)}

	writers = append(writers, zapcore.AddSync(os.Stdout))

	return zapcore.NewMultiWriteSyncer(writers...)
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
