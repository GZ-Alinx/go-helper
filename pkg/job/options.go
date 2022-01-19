package job

import (
	"context"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/logger"
	"github.com/piupuer/go-helper/pkg/utils"
)

type Options struct {
	logger          *logger.Wrapper
	ctx             context.Context
	prefix          string
	requestIdCtxKey string
	taskNameCtxKey  string
	autoRequestId   bool
}

func WithLogger(l *logger.Wrapper) func(*Options) {
	return func(options *Options) {
		if l != nil {
			getOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithCtx(ctx context.Context) func(*Options) {
	return func(options *Options) {
		if !utils.InterfaceIsNil(ctx) {
			getOptionsOrSetDefault(options).ctx = ctx
			options.logger = options.logger.WithRequestId(ctx)
		}
	}
}

func WithPrefix(prefix string) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).prefix = prefix
	}
}

func WithRequestIdCtxKey(key string) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).requestIdCtxKey = key
	}
}

func WithAutoRequestId(flag bool) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).autoRequestId = flag
	}
}

func getOptionsOrSetDefault(options *Options) *Options {
	if options == nil {
		return &Options{
			logger:          logger.NewDefaultWrapper(),
			requestIdCtxKey: constant.MiddlewareRequestIdCtxKey,
			taskNameCtxKey:  constant.JobTaskNameCtxKey,
		}
	}
	return options
}

type DriverOptions struct {
	logger *logger.Wrapper
	ctx    context.Context
	prefix string
}

func WithDriverLogger(l *logger.Wrapper) func(*DriverOptions) {
	return func(options *DriverOptions) {
		if l != nil {
			getDriverOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithDriverCtx(ctx context.Context) func(*DriverOptions) {
	return func(options *DriverOptions) {
		if !utils.InterfaceIsNil(ctx) {
			getDriverOptionsOrSetDefault(options).ctx = ctx
			options.logger = options.logger.WithRequestId(ctx)
		}
	}
}

func WithDriverPrefix(prefix string) func(*DriverOptions) {
	return func(options *DriverOptions) {
		getDriverOptionsOrSetDefault(options).prefix = prefix
	}
}

func getDriverOptionsOrSetDefault(options *DriverOptions) *DriverOptions {
	if options == nil {
		return &DriverOptions{
			logger: logger.NewDefaultWrapper(),
			prefix: constant.JobDriverPrefix,
		}
	}
	return options
}

type CronOptions struct {
	logger *logger.Wrapper
	ctx    context.Context
}

func WithCronLogger(l *logger.Wrapper) func(*CronOptions) {
	return func(options *CronOptions) {
		if l != nil {
			getCronOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithCronCtx(ctx context.Context) func(*CronOptions) {
	return func(options *CronOptions) {
		if !utils.InterfaceIsNil(ctx) {
			getCronOptionsOrSetDefault(options).ctx = ctx
			options.logger = options.logger.WithRequestId(ctx)
		}
	}
}

func getCronOptionsOrSetDefault(options *CronOptions) *CronOptions {
	if options == nil {
		return &CronOptions{
			logger: logger.NewDefaultWrapper(),
			ctx:    context.Background(),
		}
	}
	return options
}
