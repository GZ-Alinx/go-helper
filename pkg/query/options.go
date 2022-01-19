package query

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/piupuer/go-helper/ms"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/logger"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MysqlOptions struct {
	logger          *logger.Wrapper
	db              *gorm.DB
	redis           redis.UniversalClient
	cachePrefix     string
	ctx             context.Context
	enforcer        *casbin.Enforcer
	txCtxKey        string
	requestIdCtxKey string
	fsmTransition   func(ctx context.Context, logs ...resp.FsmApprovalLog) error
}

func WithMysqlLogger(l *logger.Wrapper) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		if l != nil {
			getMysqlOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithMysqlDb(db *gorm.DB) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		if db != nil {
			getMysqlOptionsOrSetDefault(options).db = db
		}
	}
}

func WithMysqlRedis(rd redis.UniversalClient) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		if rd != nil {
			getMysqlOptionsOrSetDefault(options).redis = rd
		}
	}
}

func WithMysqlCachePrefix(prefix string) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		getMysqlOptionsOrSetDefault(options).cachePrefix = prefix
	}
}

func WithMysqlCtx(ctx context.Context) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		if !utils.InterfaceIsNil(ctx) {
			getMysqlOptionsOrSetDefault(options).ctx = ctx
			options.logger = options.logger.WithRequestId(ctx)
		}
	}
}

func WithMysqlCasbinEnforcer(enforcer *casbin.Enforcer) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		if enforcer != nil {
			getMysqlOptionsOrSetDefault(options).enforcer = enforcer
		}
	}
}

func WithMysqlFsmTransition(fun func(ctx context.Context, logs ...resp.FsmApprovalLog) error) func(*MysqlOptions) {
	return func(options *MysqlOptions) {
		if fun != nil {
			getMysqlOptionsOrSetDefault(options).fsmTransition = fun
		}
	}
}

func getMysqlOptionsOrSetDefault(options *MysqlOptions) *MysqlOptions {
	if options == nil {
		return &MysqlOptions{
			logger:          logger.NewDefaultWrapper(),
			cachePrefix:     constant.QueryCachePrefix,
			ctx:             context.Background(),
			txCtxKey:        constant.MiddlewareTransactionTxCtxKey,
			requestIdCtxKey: constant.MiddlewareRequestIdCtxKey,
		}
	}
	return options
}

type MysqlReadOptions struct {
	preloads    []string
	cache       bool
	cacheExpire int
	column      string
}

func WithMySqlReadPreload(preloads ...string) func(*MysqlReadOptions) {
	return func(options *MysqlReadOptions) {
		getMysqlReadOptionsOrSetDefault(options).preloads = append(getMysqlReadOptionsOrSetDefault(options).preloads, preloads...)
	}
}

func WithMySqlReadCache(flag bool) func(*MysqlReadOptions) {
	return func(options *MysqlReadOptions) {
		getMysqlReadOptionsOrSetDefault(options).cache = flag
	}
}

func WithMySqlReadCacheExpire(seconds int) func(*MysqlReadOptions) {
	return func(options *MysqlReadOptions) {
		if seconds > 0 {
			getMysqlReadOptionsOrSetDefault(options).cacheExpire = seconds
		}
	}
}

func WithMySqlReadColumn(column string) func(*MysqlReadOptions) {
	return func(options *MysqlReadOptions) {
		if column != "" {
			getMysqlReadOptionsOrSetDefault(options).column = column
		}
	}
}

func getMysqlReadOptionsOrSetDefault(options *MysqlReadOptions) *MysqlReadOptions {
	if options == nil {
		return &MysqlReadOptions{
			preloads:    []string{},
			cacheExpire: constant.QueryCacheExpire,
			column:      constant.QueryPrimaryKey,
		}
	}
	return options
}

type RedisOptions struct {
	logger          *logger.Wrapper
	redis           redis.UniversalClient
	ctx             context.Context
	enforcer        *casbin.Enforcer
	requestIdCtxKey string
	database        string
	namingStrategy  schema.Namer
}

func WithRedisLogger(l *logger.Wrapper) func(*RedisOptions) {
	return func(options *RedisOptions) {
		if l != nil {
			getRedisOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithRedisClient(rd redis.UniversalClient) func(*RedisOptions) {
	return func(options *RedisOptions) {
		if rd != nil {
			getRedisOptionsOrSetDefault(options).redis = rd
		}
	}
}

func WithRedisCtx(ctx context.Context) func(*RedisOptions) {
	return func(options *RedisOptions) {
		if !utils.InterfaceIsNil(ctx) {
			getRedisOptionsOrSetDefault(options).ctx = ctx
			options.logger = options.logger.WithRequestId(ctx)
		}
	}
}

func WithRedisCasbinEnforcer(enforcer *casbin.Enforcer) func(*RedisOptions) {
	return func(options *RedisOptions) {
		if enforcer != nil {
			getRedisOptionsOrSetDefault(options).enforcer = enforcer
		}
	}
}

func WithRedisRequestIdCtxKey(key string) func(*RedisOptions) {
	return func(options *RedisOptions) {
		getRedisOptionsOrSetDefault(options).requestIdCtxKey = key
	}
}

func WithRedisDatabase(database string) func(*RedisOptions) {
	return func(options *RedisOptions) {
		getRedisOptionsOrSetDefault(options).database = database
	}
}

func WithRedisNamingStrategy(name schema.Namer) func(*RedisOptions) {
	return func(options *RedisOptions) {
		getRedisOptionsOrSetDefault(options).namingStrategy = name
	}
}

func getRedisOptionsOrSetDefault(options *RedisOptions) *RedisOptions {
	if options == nil {
		return &RedisOptions{
			logger:          logger.NewDefaultWrapper(),
			ctx:             context.Background(),
			requestIdCtxKey: constant.MiddlewareRequestIdCtxKey,
			database:        "query_redis",
		}
	}
	return options
}

type MessageHubOptions struct {
	logger         *logger.Wrapper
	dbNoTx         *MySql
	rd             *Redis
	idempotence    bool
	idempotenceOps []func(*middleware.IdempotenceOptions)
	findUserByIds  func(c *gin.Context, userIds []uint) []ms.User
}

func WithMessageHubLogger(l *logger.Wrapper) func(*MessageHubOptions) {
	return func(options *MessageHubOptions) {
		if l != nil {
			getMessageHubOptionsOrSetDefault(options).logger = l
		}
	}
}

func WithMessageHubDbNoTx(dbNoTx *MySql) func(*MessageHubOptions) {
	return func(options *MessageHubOptions) {
		if dbNoTx != nil {
			getMessageHubOptionsOrSetDefault(options).dbNoTx = dbNoTx
		}
	}
}

func WithMessageHubRedis(redis *Redis) func(*MessageHubOptions) {
	return func(options *MessageHubOptions) {
		if redis != nil {
			getMessageHubOptionsOrSetDefault(options).rd = redis
		}
	}
}

func WithMessageHubIdempotence(flag bool) func(*MessageHubOptions) {
	return func(options *MessageHubOptions) {
		getMessageHubOptionsOrSetDefault(options).idempotence = flag
	}
}

func WithMessageHubIdempotenceOps(ops ...func(*middleware.IdempotenceOptions)) func(*MessageHubOptions) {
	return func(options *MessageHubOptions) {
		getMessageHubOptionsOrSetDefault(options).idempotenceOps = append(getMessageHubOptionsOrSetDefault(options).idempotenceOps, ops...)
	}
}

func WithMessageHubFindUserByIds(fun func(c *gin.Context, userIds []uint) []ms.User) func(*MessageHubOptions) {
	return func(options *MessageHubOptions) {
		if fun != nil {
			getMessageHubOptionsOrSetDefault(options).findUserByIds = fun
		}
	}
}

func getMessageHubOptionsOrSetDefault(options *MessageHubOptions) *MessageHubOptions {
	if options == nil {
		return &MessageHubOptions{
			logger: logger.NewDefaultWrapper(),
		}
	}
	return options
}
