package metrics

import (
	"database/sql"
	"expvar"
	"runtime"
)

type Metrics struct {
	TotalRequests    *expvar.Int
	TotalErrors      *expvar.Int
	RequestsInFlight *expvar.Int

	PostsCreated    *expvar.Int
	UsersRegistered *expvar.Int
	CommentsCreated *expvar.Int

	ResponseTimes *expvar.Map
}

func New() *Metrics {
	// Create grouped metric maps
	httpMetrics := expvar.NewMap("http")
	businessMetrics := expvar.NewMap("business")

	m := &Metrics{
		TotalRequests:    new(expvar.Int),
		TotalErrors:      new(expvar.Int),
		RequestsInFlight: new(expvar.Int),
		PostsCreated:     new(expvar.Int),
		UsersRegistered:  new(expvar.Int),
		CommentsCreated:  new(expvar.Int),
		ResponseTimes:    new(expvar.Map).Init(),
	}

	// HTTP metrics
	httpMetrics.Set("requests_total", m.TotalRequests)
	httpMetrics.Set("errors_total", m.TotalErrors)
	httpMetrics.Set("requests_in_flight", m.RequestsInFlight)
	httpMetrics.Set("response_times_ms", m.ResponseTimes)

	// Business metrics
	businessMetrics.Set("posts_created_total", m.PostsCreated)
	businessMetrics.Set("users_registered_total", m.UsersRegistered)
	businessMetrics.Set("comments_created_total", m.CommentsCreated)

	// Runtime metrics (live data)
	expvar.Publish("runtime", expvar.Func(func() any {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		return map[string]any{
			"goroutines":        runtime.NumGoroutine(),
			"memory_alloc_mb":   memStats.Alloc / 1024 / 1024,
			"memory_sys_mb":     memStats.Sys / 1024 / 1024,
			"gc_runs":           memStats.NumGC,
			"gc_pause_total_ms": memStats.PauseTotalNs / 1_000_000,
		}
	}))

	return m
}

func (m *Metrics) RegisterDatabaseMetrics(db *sql.DB) {
	expvar.Publish("database", expvar.Func(func() any {
		stats := db.Stats()
		return map[string]any{
			"max_open":         stats.MaxOpenConnections,
			"open":             stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
			"wait_count":       stats.WaitCount,
			"wait_duration_ms": stats.WaitDuration.Milliseconds(),
		}
	}))
}
