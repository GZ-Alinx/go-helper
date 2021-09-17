package job

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/hibiken/asynq"
	"github.com/libi/dcron"
	"github.com/piupuer/go-helper/logger"
	uuid "github.com/satori/go.uuid"
	glogger "gorm.io/gorm/logger"
	"strings"
	"sync"
)

const (
	dcronInfoPrefix  = "INFO: "
	dcronErrorPrefix = "ERR: "
)

type Config struct {
	RedisUri    string
	RedisClient redis.UniversalClient
}

type GoodJob struct {
	lock   sync.Mutex
	redis  redis.UniversalClient
	driver *RedisClientDriver
	tasks  map[string]GoodTask
	ops    Options
	Error  error
}

type GoodTask struct {
	cron       *dcron.Dcron
	running    bool
	Name       string
	Expr       string
	Payload    string
	Func       func(ctx context.Context) error
	ErrHandler func(err error)
}

func New(cfg Config, options ...func(*Options)) (*GoodJob, error) {
	// init fields
	job := GoodJob{}
	if cfg.RedisClient != nil {
		job.redis = cfg.RedisClient
	} else {
		if cfg.RedisUri == "" {
			cfg.RedisUri = "redis://127.0.0.1:6379/0"
		}
		r, err := ParseRedisURI(cfg.RedisUri)
		if err != nil {
			return nil, err
		}
		job.redis = r
	}

	drv, err := NewDriver(
		job.redis,
		WithDriverLogger(job.ops.logger),
		WithDriverContext(job.ops.ctx),
		WithDriverPrefix(job.ops.prefix),
	)
	if err != nil {
		return nil, err
	}
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	job.ops = *ops
	job.driver = drv
	job.tasks = make(map[string]GoodTask, 0)
	return &job, nil
}

func ParseRedisURI(uri string) (redis.UniversalClient, error) {
	var opt asynq.RedisConnOpt
	var err error
	if uri != "" {
		opt, err = asynq.ParseRedisURI(uri)
		if err != nil {
			return nil, err
		}
		return opt.MakeRedisClient().(redis.UniversalClient), nil
	}
	return nil, fmt.Errorf("invalid redis config")
}

func (g *GoodJob) AddTask(task GoodTask) *GoodJob {
	if g.Error != nil {
		return g
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	if _, ok := g.tasks[task.Name]; ok {
		g.ops.logger.Warn(g.ops.ctx, "task %s already exists, skip", task.Name)
		return g
	}

	task.cron = dcron.NewDcronWithOption(
		task.Name,
		g.driver,
		dcron.WithLogger(&cronLogger{
			g.ops.logger,
		}),
	)
	g.tasks[task.Name] = task
	fun := (func(task GoodTask) func() {
		return func() {
			ctx := context.Background()
			if g.ops.AutoRequestId {
				ctx = context.WithValue(ctx, logger.RequestIdContextKey, uuid.NewV4().String())
			}
			err := task.Func(ctx)
			if err != nil {
				if task.ErrHandler != nil {
					task.ErrHandler(err)
				}
			}
		}
	})(task)
	task.cron.AddFunc(task.Name, task.Expr, fun)
	return g
}

func (g *GoodJob) Start() {
	if g.Error != nil {
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	for _, task := range g.tasks {
		if !task.running {
			task.cron.Start()
			task.running = true
			g.tasks[task.Name] = task
		}
	}
}

// stop all task in current node(task still running in other node)
func (g *GoodJob) StopAll() {
	if g.Error != nil {
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	for _, task := range g.tasks {
		if task.running {
			task.cron.Stop()
			task.running = false
			g.tasks[task.Name] = task
		}
	}
}

// stop task in current node(task still running in other node)
func (g *GoodJob) Stop(taskName string) {
	if g.Error != nil {
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	for _, task := range g.tasks {
		if task.running && task.Name == taskName {
			task.cron.Stop()
			task.running = false
			g.tasks[task.Name] = task
			delete(g.tasks, taskName)
			break
		} else {
			g.ops.logger.Warn(g.ops.ctx, "task %s is not running, skip", task.Name)
		}
	}
}

type cronLogger struct {
	l glogger.Interface
}

func (c cronLogger) Printf(format string, args ...interface{}) {
	ctx := context.Background()
	if strings.HasPrefix(format, dcronInfoPrefix) {
		c.l.Info(ctx, strings.TrimPrefix(format, dcronInfoPrefix), args...)
	} else if strings.HasPrefix(format, dcronErrorPrefix) {
		c.l.Error(ctx, strings.TrimPrefix(format, dcronErrorPrefix), args...)
	}
}
