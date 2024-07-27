package pgxotel

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

func NewTracer() *Tracer {
	tp := otel.GetTracerProvider()
	return &Tracer{
		tracer: tp.Tracer(TracerName, trace.WithInstrumentationVersion(InstrumentVersion)),
		attrs:  []attribute.KeyValue{semconv.DBSystemPostgreSQL},
	}
}

func recordError(span trace.Span, err error) {
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			span.SetAttributes(SQLStateKey.String(pgErr.Code))
		}
	}
}

func connAttrFromCfgPgx(config *pgx.ConnConfig) []trace.SpanStartOption {
	if config != nil {
		return []trace.SpanStartOption{
			trace.WithAttributes(
				attribute.String("database", config.Database),
				attribute.String("user", config.User),
				attribute.String("host", config.Host),
				attribute.Int("port", int(config.Port)),
			),
		}
	}

	return nil
}

func makeParamAttr(args []any) attribute.KeyValue {
	ss := make([]string, len(args))
	for i := range args {
		ss[i] = fmt.Sprintf("%+v", args[i])
	}

	return QueryParametersKey.StringSlice(ss)
}
