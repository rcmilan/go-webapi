package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const otlpEndpoint = "http://localhost:4318"

type contextKey struct{}

// Setup initialises OTel providers (traces, logs, metrics) with OTLP HTTP
// exporters pointing at the local Grafana LGTM stack, and replaces the
// default slog logger with a fan-out handler that writes JSON to stdout and
// ships log records to Loki via the OTel log bridge.
//
// Returns a shutdown function that must be called before the process exits.
func Setup(ctx context.Context) (func(), error) {
	// Route OTel internal errors through slog so they appear in the log stream.
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Error("otel", slog.String("error", err.Error()))
	}))

	// W3C TraceContext + Baggage propagation for distributed tracing.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName("bookstore-api")),
	)
	if err != nil {
		return nil, err
	}

	// ── Traces → Tempo ────────────────────────────────────────────────────────
	// WithEndpointURL forces HTTP — the SDK defaults to HTTPS.
	traceExp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(otlpEndpoint),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// ── Logs → Loki ───────────────────────────────────────────────────────────
	logExp, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpointURL(otlpEndpoint),
	)
	if err != nil {
		return nil, err
	}
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		sdklog.WithResource(res),
	)

	// ── Metrics → Prometheus ──────────────────────────────────────────────────
	metricExp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpointURL(otlpEndpoint),
	)
	if err != nil {
		return nil, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	// ── slog: JSON console + OTel bridge → Loki ───────────────────────────────
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	otelHandler := otelslog.NewHandler("bookstore-api", otelslog.WithLoggerProvider(lp))
	slog.SetDefault(slog.New(newMultiHandler(jsonHandler, otelHandler)))

	slog.Info("observability iniciada", slog.String("otlp_endpoint", otlpEndpoint))

	shutdown := func() {
		_ = tp.Shutdown(ctx)
		_ = lp.Shutdown(ctx)
		_ = mp.Shutdown(ctx)
	}
	return shutdown, nil
}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(contextKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}

// newMultiHandler fans out log records to all provided handlers.
func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
