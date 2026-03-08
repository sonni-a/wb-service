package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	KafkaMessagesProcessedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_processed_total",
			Help: "Total processed Kafka messages",
		},
	)

	KafkaProcessingErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_processing_errors_total",
			Help: "Total Kafka processing errors",
		},
	)

	KafkaDLQMessagesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_dlq_messages_total",
			Help: "Total messages sent to DLQ",
		},
	)

	CacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		},
	)

	CacheMissesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total cache misses",
		},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func Init() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
		KafkaMessagesProcessedTotal,
		KafkaProcessingErrorsTotal,
		KafkaDLQMessagesTotal,
		CacheHitsTotal,
		CacheMissesTotal,
		DBQueryDuration,
	)
}
