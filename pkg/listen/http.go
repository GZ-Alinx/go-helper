package listen

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Http(options ...func(*HttpOptions)) {
	ops := getHttpOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}

	host := ops.host
	port := ops.port
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: ops.handler,
	}

	if ops.pprofPort > 0 {
		go func() {
			// listen pprof port
			ops.logger.Info(ops.ctx, "[%s][http server]debug pprof is running at %s:%d", ops.proName, host, ops.pprofPort)
			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, ops.pprofPort), nil); err != nil {
				ops.logger.Error(ops.ctx, "[%s][http server]listen pprof error: %v", ops.proName, err)
			}
		}()
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ops.logger.Error(ops.ctx, "[%s][http server]listen error: %v", ops.proName, err)
		}
	}()

	ops.logger.Info(ops.ctx, "[%s][http server]running at %s:%d/%s", ops.proName, host, port, ops.urlPrefix)

	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ops.logger.Info(ops.ctx, "[http server]shutting down...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	_, cancel := context.WithTimeout(ops.ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ops.ctx); err != nil {
		ops.logger.Error(ops.ctx, "[%s][http server]forced to shutdown: %v", ops.proName, err)
	}

	ops.logger.Info(ops.ctx, "[%s][http server]exiting", ops.proName)
}
