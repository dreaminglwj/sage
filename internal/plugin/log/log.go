package log

import (
	"fmt"
	"time"

	"github.com/dreaminglwj/sage/internal/resource"
	"github.com/dreaminglwj/sage/internal/tracing"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	opentracing "github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	xormLog "xorm.io/xorm/log"

	"github.com/dreaminglwj/sage/internal/conf"
)

var _ kratosLog.Logger = (*Logger)(nil)
var _ xormLog.ContextLogger = (*Logger)(nil)

type Logger struct {
	*zap.SugaredLogger
	helper    *kratosLog.Helper
	level     xormLog.LogLevel
	isShowSql bool
}

func (l *Logger) BeforeSQL(ctx xormLog.LogContext) {
	_, c := opentracing.StartSpanFromContext(ctx.Ctx, "XORM SQL execute")
	ctx.Ctx = c
}

func (l *Logger) AfterSQL(ctx xormLog.LogContext) {
	span := opentracing.SpanFromContext(ctx.Ctx)
	if span == nil {
		if ctx.ExecuteTime > 0 {
			l.helper.WithContext(ctx.Ctx).Infof("[SQL]%s %v - %v", ctx.SQL, ctx.Args, ctx.ExecuteTime)
		} else {
			l.helper.WithContext(ctx.Ctx).Infof("[SQL]%s %v", ctx.SQL, ctx.Args)
		}
		return
	}
	defer span.Finish()
	var sessionPart string
	v := ctx.Ctx.Value("__xorm_session_id")
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
		span.LogFields(tracerLog.String("session_id", sessionPart))
	}

	span.LogFields(tracerLog.String("SQL", ctx.SQL))
	span.LogFields(tracerLog.Object("args", ctx.Args))
	span.SetTag("execute_time", ctx.ExecuteTime)

	if ctx.ExecuteTime > 0 {
		l.helper.WithContext(ctx.Ctx).Infof("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.helper.WithContext(ctx.Ctx).Infof("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	}
}

func (l *Logger) Close() error {
	_ = l.Sync()
	return nil
}

// NewLogger Logger constructor
func NewLogger(config *conf.Config) (*Logger, error) {
	var (
		zapConfig zap.Config
	)
	options := []zap.Option{
		zap.AddCallerSkip(3),
	}
	if config.App.IsProduction() {
		zapConfig = zap.NewProductionConfig()
		options = append(options, zap.Fields(
			zap.String("_VER_", config.App.Version),
			zap.String("_DEPLOY_", string(config.App.Env)),
			zap.String("_APP_", config.App.Name),
			zap.String("_COMP_", config.App.Component),
			zap.String("_TYPE_", config.Log.Type),
		))
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	zapConfig.EncoderConfig.LevelKey = "_LEVEL_"
	zapConfig.EncoderConfig.TimeKey = "_TS_"
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	zapConfig.EncoderConfig.NameKey = "_NAME_"
	zapConfig.EncoderConfig.MessageKey = "_MSG_"
	zapConfig.EncoderConfig.CallerKey = "_CALLER_"
	zapConfig.EncoderConfig.StacktraceKey = "_STACKTRACE_"
	logger, err := zapConfig.Build(options...)
	if err != nil {
		return nil, err
	}

	l := &Logger{SugaredLogger: logger.Sugar()}
	l.helper = kratosLog.NewHelper(kratosLog.With(l,
		"traceID", tracing.TraceID(),
		"spanID", tracing.SpanID(),
	))
	resource.Register(l)
	return l, nil
}

func NewHelper(logger *Logger) *kratosLog.Helper {
	return logger.helper
}

// Level implement xorm.io/xorm/log.Logger
func (l *Logger) Level() xormLog.LogLevel {
	return l.level
}

// SetLevel implement xorm.io/xorm/log.Logger
func (l *Logger) SetLevel(level xormLog.LogLevel) {
	l.level = level
}

// ShowSQL implement xorm.io/xorm/log.Logger
func (l *Logger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.isShowSql = true
		return
	}
	l.isShowSql = show[0]
}

// IsShowSQL implement xorm.io/xorm/log.Logger
func (l *Logger) IsShowSQL() bool {
	return l.isShowSql
}

// Log implement github.com/go-kratos/kratos/v2/log.Logger
func (l *Logger) Log(level kratosLog.Level, keyvals ...any) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	switch level {
	case kratosLog.LevelDebug:
		l.Desugar().Debug("", data...)
	case kratosLog.LevelInfo:
		l.Desugar().Info("", data...)
	case kratosLog.LevelWarn:
		l.Desugar().Warn("", data...)
	case kratosLog.LevelError:
		l.Desugar().Error("", data...)
	case kratosLog.LevelFatal:
		l.Desugar().Fatal("", data...)
	}
	return nil
}
