package mq

import (
	"github.com/streadway/amqp"
	"github.com/thoas/go-funk"
)

type RabbitOptions struct {
	QosPrefetchCount       int
	QosPrefetchSize        int
	QosGlobal              bool
	ReconnectInterval      int
	ReconnectMaxRetryCount int
	Timeout                int
}

func WithChannelQosPrefetchCount(prefetchCount int) func(*RabbitOptions) {
	return func(options *RabbitOptions) {
		getRabbitOptionsOrSetDefault(options).QosPrefetchCount = prefetchCount
	}
}

func WithChannelQosPrefetchSize(prefetchSize int) func(*RabbitOptions) {
	return func(options *RabbitOptions) {
		getRabbitOptionsOrSetDefault(options).QosPrefetchSize = prefetchSize
	}
}

func WithChannelQosGlobal(options *RabbitOptions) {
	getRabbitOptionsOrSetDefault(options).QosGlobal = true
}

func WithReconnectInterval(second int) func(*RabbitOptions) {
	return func(options *RabbitOptions) {
		if second > 0 {
			getRabbitOptionsOrSetDefault(options).ReconnectInterval = second
		}
	}
}

func WithReconnectMaxRetryCount(count int) func(*RabbitOptions) {
	return func(options *RabbitOptions) {
		if count > 0 {
			getRabbitOptionsOrSetDefault(options).ReconnectMaxRetryCount = count
		}
	}
}

func WithTimeout(second int) func(*RabbitOptions) {
	return func(options *RabbitOptions) {
		if second > 0 {
			getRabbitOptionsOrSetDefault(options).Timeout = second
		}
	}
}

func getRabbitOptionsOrSetDefault(options *RabbitOptions) *RabbitOptions {
	if options == nil {
		return &RabbitOptions{
			QosPrefetchCount:       2,
			Timeout:                10,
			ReconnectMaxRetryCount: 1,
			ReconnectInterval:      5,
		}
	}
	return options
}

type ExchangeOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
	Declare    bool
	NamePrefix string
}

func WithExchangeName(name string) func(*ExchangeOptions) {
	return func(options *ExchangeOptions) {
		getExchangeOptionsOrSetDefault(options).Name = name
	}
}

func WithExchangeKind(kind string) func(*ExchangeOptions) {
	return func(options *ExchangeOptions) {
		getExchangeOptionsOrSetDefault(options).Kind = kind
	}
}

func WithExchangeDurable(options *ExchangeOptions) {
	getExchangeOptionsOrSetDefault(options).Durable = true
}

func WithExchangeAutoDelete(options *ExchangeOptions) {
	getExchangeOptionsOrSetDefault(options).AutoDelete = true
}

func WithExchangeInternal(options *ExchangeOptions) {
	getExchangeOptionsOrSetDefault(options).Internal = true
}

func WithExchangeNoWait(options *ExchangeOptions) {
	getExchangeOptionsOrSetDefault(options).NoWait = true
}

func WithExchangeArgs(args amqp.Table) func(*ExchangeOptions) {
	return func(options *ExchangeOptions) {
		getExchangeOptionsOrSetDefault(options).Args = args
	}
}

func WithExchangeSkipDeclare(options *ExchangeOptions) {
	getExchangeOptionsOrSetDefault(options).Declare = false
}

func WithExchangeNamePrefix(prefix string) func(*ExchangeOptions) {
	return func(options *ExchangeOptions) {
		getExchangeOptionsOrSetDefault(options).NamePrefix = prefix
	}
}

func getExchangeOptionsOrSetDefault(options *ExchangeOptions) *ExchangeOptions {
	if options == nil {
		return &ExchangeOptions{
			Kind:    amqp.ExchangeDirect,
			Durable: true,
			Declare: true,
		}
	}
	return options
}

type QueueOptions struct {
	Name           string
	RouteKeys      []string
	Durable        bool
	AutoDelete     bool
	Exclusive      bool
	NoWait         bool
	Args           amqp.Table
	BindArgs       amqp.Table
	Declare        bool
	Bind           bool
	DeadLetterName string
	DeadLetterKey  string
	MessageTTL     int32
	NamePrefix     string
}

func WithQueueName(name string) func(*QueueOptions) {
	return func(options *QueueOptions) {
		getQueueOptionsOrSetDefault(options).Name = name
	}
}

func WithQueueRouteKey(key string) func(*QueueOptions) {
	return func(options *QueueOptions) {
		d := getQueueOptionsOrSetDefault(options)
		keys := d.RouteKeys
		if !funk.ContainsString(keys, key) {
			d.RouteKeys = append(keys, key)
		}
	}
}

func WithQueueSkipDeclare(options *QueueOptions) {
	getQueueOptionsOrSetDefault(options).Declare = false
}

func WithQueueSkipBind(options *QueueOptions) {
	getQueueOptionsOrSetDefault(options).Bind = false
}

func WithQueueDeadLetterName(name string) func(*QueueOptions) {
	return func(options *QueueOptions) {
		getQueueOptionsOrSetDefault(options).DeadLetterName = name
	}
}

func WithQueueDeadLetterKey(key string) func(*QueueOptions) {
	return func(options *QueueOptions) {
		getQueueOptionsOrSetDefault(options).DeadLetterKey = key
	}
}

func WithQueueMessageTTL(ttl int32) func(*QueueOptions) {
	return func(options *QueueOptions) {
		getQueueOptionsOrSetDefault(options).MessageTTL = ttl
	}
}

func getQueueOptionsOrSetDefault(options *QueueOptions) *QueueOptions {
	if options == nil {
		return &QueueOptions{
			Durable: true,
			Declare: true,
			Bind:    true,
		}
	}
	return options
}

type PublishOptions struct {
	RouteKeys    []string
	ContentType  string
	Headers      amqp.Table
	DeliveryMode uint8
	Mandatory    bool
	Immediate    bool
	Expiration   string
}

func WithPublishOptionsContentType(contentType string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		getPublishOptionsOrSetDefault(options).ContentType = contentType
	}
}

func WithPublishOptionsHeaders(headers amqp.Table) func(*PublishOptions) {
	return func(options *PublishOptions) {
		getPublishOptionsOrSetDefault(options).Headers = headers
	}
}

func WithPublishRouteKey(key string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		d := getPublishOptionsOrSetDefault(options)
		keys := d.RouteKeys
		if !funk.ContainsString(keys, key) {
			d.RouteKeys = append(keys, key)
		}
	}
}

func getPublishOptionsOrSetDefault(options *PublishOptions) *PublishOptions {
	if options == nil {
		return &PublishOptions{
			ContentType: "text/plain",
		}
	}
	return options
}
