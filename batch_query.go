package pgxotel

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

func (t *Tracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	var size int
	if b := data.Batch; b != nil {
		size = b.Len()
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(BatchSizeKey.Int(size)),
	}

	if conn != nil {
		opts = append(opts, connAttrFromCfgPgx(conn.Config())...)
	}

	ctx, _ = t.tracer.Start(ctx, "BATCH QUERY START ", opts...)

	return ctx
}

func (t *Tracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}

	if conn != nil {
		opts = append(opts, connAttrFromCfgPgx(conn.Config())...)
		opts = append(opts, trace.WithAttributes(DBStatement.String(data.SQL)))
		opts = append(opts, trace.WithAttributes(makeParamAttr(data.Args)))
	}

	spanName := "BATCH QUERY | " + data.SQL

	_, span := t.tracer.Start(ctx, spanName, opts...)
	recordError(span, data.Err)

	span.End()
}

func (t *Tracer) TraceBatchEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceBatchEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)

	span.End()
}
