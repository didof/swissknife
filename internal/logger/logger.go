/*
<https://www.golinuxcloud.com/golang-zap-logger/>
<https://github.com/bilalcaliskan/syn-flood/blob/master/internal/logging/logging.go>
*/

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	atomic zap.AtomicLevel
)

func init() {
	atomic = zap.NewAtomicLevel()
	atomic.SetLevel(zap.InfoLevel)
	logger = zap.New(zapcore.NewTee(zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "severity",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.RFC3339TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
	}), zapcore.Lock(os.Stdout), atomic)))
}

func SetLevel(l zapcore.Level) {
	atomic.SetLevel(l)
}

func Get() *zap.Logger {
	return logger
}
