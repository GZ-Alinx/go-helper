package oss

import (
	"context"
	"github.com/piupuer/go-helper/pkg/logger"
)

type MinioOptions struct {
	logger   logger.Interface
	ctx      context.Context
	endpoint string
	accessId string
	secret   string
	https    bool
}

func WithMinioLogger(l logger.Interface) func(*MinioOptions) {
	return func(options *MinioOptions) {
		if l != nil {
			getMinioOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithMinioLoggerLevel(level logger.Level) func(*MinioOptions) {
	return func(options *MinioOptions) {
		l := options.logger
		if options.logger == nil {
			l = getMinioOptionsOrSetDefault(options).logger
		}
		options.logger = l.LogLevel(level)
	}
}

func WithMinioContext(ctx context.Context) func(*MinioOptions) {
	return func(options *MinioOptions) {
		getMinioOptionsOrSetDefault(options).ctx = ctx
	}
}

func WithMinioEndpoint(endpoint string) func(*MinioOptions) {
	return func(options *MinioOptions) {
		getMinioOptionsOrSetDefault(options).endpoint = endpoint
	}
}

func WithMinioAccessId(accessId string) func(*MinioOptions) {
	return func(options *MinioOptions) {
		getMinioOptionsOrSetDefault(options).accessId = accessId
	}
}

func WithMinioSecret(secret string) func(*MinioOptions) {
	return func(options *MinioOptions) {
		getMinioOptionsOrSetDefault(options).secret = secret
	}
}

func WithMinioHttps(options *MinioOptions) {
	getMinioOptionsOrSetDefault(options).https = true
}

func getMinioOptionsOrSetDefault(options *MinioOptions) *MinioOptions {
	if options == nil {
		return &MinioOptions{
			logger: logger.DefaultLogger(),
		}
	}
	return options
}
