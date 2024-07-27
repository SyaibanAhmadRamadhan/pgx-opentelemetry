package pgxotel

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

func (t *Tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(DBStatement.String(data.SQL)),
		trace.WithAttributes(makeParamAttr(data.Args)),
	}

	if conn != nil {
		opts = append(opts, connAttrFromCfgPgx(conn.Config())...)
	}

	spanName := "QUERY | " + data.SQL
	ctx, _ = t.tracer.Start(ctx, spanName, opts...)

	return ctx
}

func (t *Tracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)

	if data.Err != nil {
		span.SetAttributes(RowsAffectedKey.Int64(data.CommandTag.RowsAffected()))
	}

	span.End()
}
