package middleware

import (
	"github.com/piupuer/go-helper/pkg/tracing"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer(tracing.Middleware)
