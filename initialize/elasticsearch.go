package initialize

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/supuwoerc/weaver/conf"
	weaverLogger "github.com/supuwoerc/weaver/pkg/logger"
)

type ElasticsearchLogLevel int

const (
	None ElasticsearchLogLevel = iota
	Request
	Response
	All
)

type ElasticsearchLogger struct {
	goredislib.Hook
	*weaverLogger.Logger
	Level ElasticsearchLogLevel
}

func NewElasticsearchLogger(l *weaverLogger.Logger, conf *conf.Config) *ElasticsearchLogger {
	return &ElasticsearchLogger{
		Logger: l,
		Level:  ElasticsearchLogLevel(conf.Elasticsearch.LogLevel),
	}
}

func (l *ElasticsearchLogger) RequestBodyEnabled() bool { return l.Level == Request || l.Level == All }
func (l *ElasticsearchLogger) ResponseBodyEnabled() bool {
	return l.Level == Response || l.Level == All
}

func (l *ElasticsearchLogger) LogRoundTrip(req *http.Request, res *http.Response, err error, _ time.Time, dur time.Duration) error {
	if err != nil {
		l.Errorw("elasticsearch request failed", "err", err)
	} else {
		l.Infow("elasticsearch request success", "url", req.URL, "duration", dur, "status_code", res.StatusCode)
	}
	return nil
}

func NewElasticsearchClient(conf *conf.Config, logger elastictransport.Logger) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses:           conf.Elasticsearch.Addresses,
		Username:            conf.Elasticsearch.Username,
		Password:            conf.Elasticsearch.Password,
		APIKey:              conf.Elasticsearch.APIKey,
		ServiceToken:        conf.Elasticsearch.ServiceToken,
		MaxRetries:          conf.Elasticsearch.MaxRetries,
		RetryOnStatus:       []int{502, 503, 504, 429},
		CompressRequestBody: conf.Elasticsearch.CompressRequestBody,
		EnableMetrics:       conf.Elasticsearch.EnableMetrics,
		EnableDebugLogger:   conf.Elasticsearch.EnableDebugLogger,
		Logger:              logger,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: conf.Elasticsearch.Insecure,
			},
		},
	}
	if conf.Elasticsearch.DiscoverNodesOnStart {
		cfg.EnableDebugLogger = true
		if conf.Elasticsearch.DiscoverNodesInterval > 0 {
			cfg.DiscoverNodesInterval = conf.Elasticsearch.DiscoverNodesInterval * time.Second
		}
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return client
}
