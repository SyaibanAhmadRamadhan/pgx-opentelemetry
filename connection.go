package pgxotel

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

func (t *Tracer) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}

	if data.ConnConfig != nil {
		opts = append(opts, connAttrFromCfgPgx(data.ConnConfig)...)
	}

	ctx, _ = t.tracer.Start(ctx, "CONNECT ", opts...)

	return ctx
}

func (t *Tracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)

	span.End()
}
