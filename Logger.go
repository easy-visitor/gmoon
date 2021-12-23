package gmoon

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLogger struct {
	ZapConfig ZapConfig
}

func NewZapLogger(config ZapConfig) *ZapLogger {
	return &ZapLogger{
		ZapConfig: config,
	}
}

func (this *ZapLogger) Logger() (logger *zap.Logger) {

	if ok, _ := PathExists(this.ZapConfig.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", this.ZapConfig.Director)
		_ = os.Mkdir(this.ZapConfig.Director, os.ModePerm)
	}
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{
		this.getEncoderCore(fmt.Sprintf("./%s/server_debug.log", this.ZapConfig.Director), debugPriority),
		this.getEncoderCore(fmt.Sprintf("./%s/server_info.log", this.ZapConfig.Director), infoPriority),
		this.getEncoderCore(fmt.Sprintf("./%s/server_warn.log", this.ZapConfig.Director), warnPriority),
		this.getEncoderCore(fmt.Sprintf("./%s/server_error.log", this.ZapConfig.Director), errorPriority),
	}
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	if this.ZapConfig.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func (this *ZapLogger) getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  this.ZapConfig.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     this.customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case this.ZapConfig.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case this.ZapConfig.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case this.ZapConfig.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case this.ZapConfig.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func (this *ZapLogger) getEncoder() zapcore.Encoder {
	if this.ZapConfig.Format == "json" {
		return zapcore.NewJSONEncoder(this.getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(this.getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func (this *ZapLogger) getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := this.getWriteSyncer(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(this.getEncoder(), writer, level)
}

func (this *ZapLogger) getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件的位置
		MaxSize:    10,   // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  // 保留旧文件的最大个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
	}

	if this.ZapConfig.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 自定义日志输出时间格式
func (this *ZapLogger) customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(this.ZapConfig.Prefix + "2006/01/02 - 15:04:05.000"))
}
