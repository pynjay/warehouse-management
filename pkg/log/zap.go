package log

import (
	"fmt"
	"sync"
	"time"
	"warehouse/pkg/errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*ZapWrapper)(nil)

var DefaultLogger, _ = NewZapWrapper("", false)

type ZapWrapper struct {
	logger   *zap.SugaredLogger
	core     *zapCoreWithCustomStack
	cloneMtx sync.Mutex
}

func NewZapWrapper(dst string, debug bool) (*ZapWrapper, error) {
	if dst == "" {
		dst = "stdout"
	}

	var loggerConfig zap.Config
	if debug {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	loggerConfig.DisableStacktrace = true

	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.OutputPaths = []string{dst}
	loggerConfig.Development = debug

	var core *zapCoreWithCustomStack

	logger, err := loggerConfig.Build(
		zap.AddCallerSkip(1),
		zap.WrapCore(
			func(c zapcore.Core) zapcore.Core {
				core = &zapCoreWithCustomStack{core: c}
				return core
			},
		),
	)
	if err != nil {
		return nil, err
	}

	lg := logger.Sugar()

	return &ZapWrapper{logger: lg, core: core}, nil
}

func (z *ZapWrapper) Warn(kv ...interface{}) {
	if len(kv) == 0 {
		z.logger.Warn(kv...)
		return
	}

	if len(kv)%2 != 0 {
		switch val := kv[0].(type) {
		case error:
			_, stacker := errors.DeepestErrorWithStack(val)
			if stacker != nil {
				z.core.SetStackAndLock([]byte(fmt.Sprintf("%+v", stacker)))
				defer z.core.UnlockStack()
				// kv = append(kv, "err_stacktrace", fmt.Sprintf("%+v", stacker))
			}
			z.logger.Warnw(val.Error(), kv[1:]...)
			return
		case fmt.Stringer:
			z.logger.Warnw(val.String(), kv[1:]...)
			return

		case string:
			z.logger.Warnw(val, kv[1:]...)
			return
		}
	}

	z.logger.Warn(kv...)
}

func (z *ZapWrapper) Err(kv ...interface{}) {
	if len(kv) == 0 {
		z.logger.Error(kv...)
		return
	}

	if len(kv)%2 != 0 {
		switch val := kv[0].(type) {
		case error:
			_, stacker := errors.DeepestErrorWithStack(val)
			if stacker != nil {
				z.core.SetStackAndLock([]byte(fmt.Sprintf("%+v", stacker)))
				defer z.core.UnlockStack()
				// kv = append(kv, "err_stacktrace", fmt.Sprintf("%+v", stacker))
			}
			z.logger.Errorw(val.Error(), kv[1:]...)
			return

		case fmt.Stringer:
			z.logger.Errorw(val.String(), kv[1:]...)
			return

		case string:
			z.logger.Errorw(val, kv[1:]...)
			return
		}
	}

	z.logger.Error(kv...)
}

func (z *ZapWrapper) Debug(kv ...interface{}) {
	if len(kv) == 0 {
		z.logger.Debug(kv...)
		return
	}

	if len(kv)%2 != 0 {
		if msg, isStr := kv[0].(string); isStr {
			z.logger.Debugw(msg, kv[1:]...)
			return
		}
	}

	z.logger.Debug(kv...)
}

func (z *ZapWrapper) Info(kv ...interface{}) {
	if len(kv) == 0 {
		z.logger.Info(kv...)
		return
	}

	if len(kv)%2 != 0 {
		if msg, isStr := kv[0].(string); isStr {
			z.logger.Infow(msg, kv[1:]...)
			return
		}
	}

	z.logger.Info(kv...)
}

func (z *ZapWrapper) Printf(format string, v ...interface{}) {
	z.logger.Infof(format, v...)
}

func (z *ZapWrapper) WithPrefix(v ...string) Logger {
	vi := make([]interface{}, 0, len(v))
	for _, i := range v {
		vi = append(vi, i)
	}

	z.cloneMtx.Lock()
	defer z.cloneMtx.Unlock()

	logger := ZapWrapper{
		logger:   z.logger.With(vi...),
		cloneMtx: sync.Mutex{},
	}

	logger.core = z.core.cloned

	return &logger
}
