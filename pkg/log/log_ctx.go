package log

import (
	"context"
	"go.uber.org/zap"
)

type loggerKeyType int

const loggerKey = iota

func Ctx(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return zap.S()
	}
	if loggerFromCtx, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return loggerFromCtx
	}
	return zap.S()
}

func WithCtx(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
